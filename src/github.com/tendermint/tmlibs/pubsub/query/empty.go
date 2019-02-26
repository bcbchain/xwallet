package query

import "github.com/tendermint/tmlibs/pubsub"

type Empty struct {
}

func (Empty) Matches(tags pubsub.TagMap) bool {
	return true
}

func (Empty) String() string {
	return "empty"
}
