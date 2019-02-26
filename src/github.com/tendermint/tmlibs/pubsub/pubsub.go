package pubsub

import (
	"context"
	"errors"
	"sync"

	cmn "github.com/tendermint/tmlibs/common"
)

type operation int

const (
	sub	operation	= iota
	pub
	unsub
	shutdown
)

var (
	ErrSubscriptionNotFound	= errors.New("subscription not found")

	ErrAlreadySubscribed	= errors.New("already subscribed")
)

type TagMap interface {
	Get(key string) (value interface{}, ok bool)

	Len() int
}

type tagMap map[string]interface{}

type cmd struct {
	op		operation
	query		Query
	ch		chan<- interface{}
	clientID	string
	msg		interface{}
	tags		TagMap
}

type Query interface {
	Matches(tags TagMap) bool
	String() string
}

type Server struct {
	cmn.BaseService

	cmds	chan cmd
	cmdsCap	int

	mtx		sync.RWMutex
	subscriptions	map[string]map[string]Query
}

type Option func(*Server)

func NewTagMap(data map[string]interface{}) TagMap {
	return tagMap(data)
}

func (ts tagMap) Get(key string) (value interface{}, ok bool) {
	value, ok = ts[key]
	return
}

func (ts tagMap) Len() int {
	return len(ts)
}

func NewServer(options ...Option) *Server {
	s := &Server{
		subscriptions: make(map[string]map[string]Query),
	}
	s.BaseService = *cmn.NewBaseService(nil, "PubSub", s)

	for _, option := range options {
		option(s)
	}

	s.cmds = make(chan cmd, s.cmdsCap)

	return s
}

func BufferCapacity(cap int) Option {
	return func(s *Server) {
		if cap > 0 {
			s.cmdsCap = cap
		}
	}
}

func (s *Server) BufferCapacity() int {
	return s.cmdsCap
}

func (s *Server) Subscribe(ctx context.Context, clientID string, query Query, out chan<- interface{}) error {
	s.mtx.RLock()
	clientSubscriptions, ok := s.subscriptions[clientID]
	if ok {
		_, ok = clientSubscriptions[query.String()]
	}
	s.mtx.RUnlock()
	if ok {
		return ErrAlreadySubscribed
	}

	select {
	case s.cmds <- cmd{op: sub, clientID: clientID, query: query, ch: out}:
		s.mtx.Lock()
		if _, ok = s.subscriptions[clientID]; !ok {
			s.subscriptions[clientID] = make(map[string]Query)
		}
		s.subscriptions[clientID][query.String()] = query
		s.mtx.Unlock()
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *Server) Unsubscribe(ctx context.Context, clientID string, query Query) error {
	var origQuery Query
	s.mtx.RLock()
	clientSubscriptions, ok := s.subscriptions[clientID]
	if ok {
		origQuery, ok = clientSubscriptions[query.String()]
	}
	s.mtx.RUnlock()
	if !ok {
		return ErrSubscriptionNotFound
	}

	select {
	case s.cmds <- cmd{op: unsub, clientID: clientID, query: origQuery}:
		s.mtx.Lock()
		delete(clientSubscriptions, query.String())
		s.mtx.Unlock()
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *Server) UnsubscribeAll(ctx context.Context, clientID string) error {
	s.mtx.RLock()
	_, ok := s.subscriptions[clientID]
	s.mtx.RUnlock()
	if !ok {
		return ErrSubscriptionNotFound
	}

	select {
	case s.cmds <- cmd{op: unsub, clientID: clientID}:
		s.mtx.Lock()
		delete(s.subscriptions, clientID)
		s.mtx.Unlock()
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *Server) Publish(ctx context.Context, msg interface{}) error {
	return s.PublishWithTags(ctx, msg, NewTagMap(make(map[string]interface{})))
}

func (s *Server) PublishWithTags(ctx context.Context, msg interface{}, tags TagMap) error {
	select {
	case s.cmds <- cmd{op: pub, msg: msg, tags: tags}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *Server) OnStop() {
	s.cmds <- cmd{op: shutdown}
}

type state struct {
	queries	map[Query]map[string]chan<- interface{}

	clients	map[string]map[Query]struct{}
}

func (s *Server) OnStart() error {
	go s.loop(state{
		queries:	make(map[Query]map[string]chan<- interface{}),
		clients:	make(map[string]map[Query]struct{}),
	})
	return nil
}

func (s *Server) OnReset() error {
	return nil
}

func (s *Server) loop(state state) {
loop:
	for cmd := range s.cmds {
		switch cmd.op {
		case unsub:
			if cmd.query != nil {
				state.remove(cmd.clientID, cmd.query)
			} else {
				state.removeAll(cmd.clientID)
			}
		case shutdown:
			for clientID := range state.clients {
				state.removeAll(clientID)
			}
			break loop
		case sub:
			state.add(cmd.clientID, cmd.query, cmd.ch)
		case pub:
			state.send(cmd.msg, cmd.tags)
		}
	}
}

func (state *state) add(clientID string, q Query, ch chan<- interface{}) {

	if _, ok := state.queries[q]; !ok {
		state.queries[q] = make(map[string]chan<- interface{})
	}

	state.queries[q][clientID] = ch

	if _, ok := state.clients[clientID]; !ok {
		state.clients[clientID] = make(map[Query]struct{})
	}
	state.clients[clientID][q] = struct{}{}
}

func (state *state) remove(clientID string, q Query) {
	clientToChannelMap, ok := state.queries[q]
	if !ok {
		return
	}

	ch, ok := clientToChannelMap[clientID]
	if ok {
		close(ch)

		delete(state.clients[clientID], q)

		if len(state.clients[clientID]) == 0 {
			delete(state.clients, clientID)
		}

		delete(state.queries[q], clientID)
		if len(state.queries[q]) == 0 {
			delete(state.queries, q)
		}
	}
}

func (state *state) removeAll(clientID string) {
	queryMap, ok := state.clients[clientID]
	if !ok {
		return
	}

	for q := range queryMap {
		ch := state.queries[q][clientID]
		close(ch)

		delete(state.queries[q], clientID)
		if len(state.queries[q]) == 0 {
			delete(state.queries, q)
		}
	}

	delete(state.clients, clientID)
}

func (state *state) send(msg interface{}, tags TagMap) {
	for q, clientToChannelMap := range state.queries {
		if q.Matches(tags) {
			for _, ch := range clientToChannelMap {
				ch <- msg
			}
		}
	}
}
