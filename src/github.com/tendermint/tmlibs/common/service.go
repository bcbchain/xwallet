package common

import (
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/tendermint/tmlibs/log"
)

var (
	ErrAlreadyStarted	= errors.New("already started")
	ErrAlreadyStopped	= errors.New("already stopped")
)

type Service interface {
	Start() error
	OnStart() error

	Stop() error
	OnStop()

	Reset() error
	OnReset() error

	IsRunning() bool

	Quit() <-chan struct{}

	String() string

	SetLogger(log.Logger)
}

type BaseService struct {
	Logger	log.Logger
	name	string
	started	uint32
	stopped	uint32
	quit	chan struct{}

	impl	Service
}

func NewBaseService(logger log.Logger, name string, impl Service) *BaseService {
	if logger == nil {
		logger = log.NewNopLogger()
	}

	return &BaseService{
		Logger:	logger,
		name:	name,
		quit:	make(chan struct{}),
		impl:	impl,
	}
}

func (bs *BaseService) SetLogger(l log.Logger) {
	bs.Logger = l
}

func (bs *BaseService) Start() error {
	if atomic.CompareAndSwapUint32(&bs.started, 0, 1) {
		if atomic.LoadUint32(&bs.stopped) == 1 {
			bs.Logger.Error(Fmt("Not starting %v -- already stopped", bs.name), "impl", bs.impl)
			return ErrAlreadyStopped
		}
		bs.Logger.Info(Fmt("Starting %v", bs.name), "impl", bs.impl)
		err := bs.impl.OnStart()
		if err != nil {

			atomic.StoreUint32(&bs.started, 0)
			return err
		}
		return nil
	}
	bs.Logger.Debug(Fmt("Not starting %v -- already started", bs.name), "impl", bs.impl)
	return ErrAlreadyStarted
}

func (bs *BaseService) OnStart() error	{ return nil }

func (bs *BaseService) Stop() error {
	if atomic.CompareAndSwapUint32(&bs.stopped, 0, 1) {
		bs.Logger.Info(Fmt("Stopping %v", bs.name), "impl", bs.impl)
		bs.impl.OnStop()
		close(bs.quit)
		return nil
	}
	bs.Logger.Debug(Fmt("Stopping %v (ignoring: already stopped)", bs.name), "impl", bs.impl)
	return ErrAlreadyStopped
}

func (bs *BaseService) OnStop()	{}

func (bs *BaseService) Reset() error {
	if !atomic.CompareAndSwapUint32(&bs.stopped, 1, 0) {
		bs.Logger.Debug(Fmt("Can't reset %v. Not stopped", bs.name), "impl", bs.impl)
		return fmt.Errorf("can't reset running %s", bs.name)
	}

	atomic.CompareAndSwapUint32(&bs.started, 1, 0)

	bs.quit = make(chan struct{})
	return bs.impl.OnReset()
}

func (bs *BaseService) OnReset() error {
	PanicSanity("The service cannot be reset")
	return nil
}

func (bs *BaseService) IsRunning() bool {
	return atomic.LoadUint32(&bs.started) == 1 && atomic.LoadUint32(&bs.stopped) == 0
}

func (bs *BaseService) Wait() {
	<-bs.quit
}

func (bs *BaseService) String() string {
	return bs.name
}

func (bs *BaseService) Quit() <-chan struct{} {
	return bs.quit
}
