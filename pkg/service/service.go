package service

import (
	"errors"
	"fmt"
	"github.com/aioncore/ouroboros/pkg/service/log"
	"sync/atomic"
)

// Service defines a handlers in a node, everything in ouroboros expect core is a handlers
type Service interface {
	// Start the handlers.
	// If it's already started or stopped, will return an error.
	// If OnStart() returns an error, it's returned by Start()
	Start() error
	OnStart() error

	// Stop the handlers.
	// If it's already stopped, will return an error.
	// OnStop must never error.
	Stop() error
	OnStop()

	// Reset the handlers.
	// Panics by default - must be overwritten to enable reset.
	Reset() error
	OnReset() error

	// IsRunning Return true if the handlers is running
	IsRunning() bool

	// Quit returns a channel, which is closed once handlers is stopped.
	Quit() <-chan struct{}

	// String representation of the handlers
	String() string
}

type BaseService struct {
	name    string
	started uint32 // atomic
	stopped uint32 // atomic
	quit    chan struct{}

	// The "subclass" of BaseService
	impl Service
}

// NewBaseService creates a new BaseService.
func NewBaseService(name string, impl Service) *BaseService {
	return &BaseService{
		name: name,
		quit: make(chan struct{}),
		impl: impl,
	}
}

// Start implements Service by calling OnStart (if defined). An error will be
// returned if the handlers is already running or stopped. Not to start the
// stopped handlers, you need to call Reset.
func (bs *BaseService) Start() error {
	if atomic.CompareAndSwapUint32(&bs.started, 0, 1) {
		if atomic.LoadUint32(&bs.stopped) == 1 {
			log.Error(fmt.Sprintf("Not starting %v handlers -- already stopped", bs.name))
			// revert flag
			atomic.StoreUint32(&bs.started, 0)
			return errors.New("already stopped")
		}
		log.Info("handlers start")
		err := bs.impl.OnStart()
		if err != nil {
			// revert flag
			atomic.StoreUint32(&bs.started, 0)
			return err
		}
		return nil
	}
	log.Debug("handlers start")
	return errors.New("already started")
}

// OnStart implements Service by doing nothing.
// NOTE: Do not put anything in here,
// that way users don't need to call BaseService.OnStart()
func (bs *BaseService) OnStart() error { return nil }

// Stop implements Service by calling OnStop (if defined) and closing quit
// channel. An error will be returned if the handlers is already stopped.
func (bs *BaseService) Stop() error {
	if atomic.CompareAndSwapUint32(&bs.stopped, 0, 1) {
		if atomic.LoadUint32(&bs.started) == 0 {
			log.Error(fmt.Sprintf("Not stopping %v handlers -- has not been started yet", bs.name))
			// revert flag
			atomic.StoreUint32(&bs.stopped, 0)
			return errors.New("not started")
		}
		log.Info("handlers stop")
		bs.impl.OnStop()
		close(bs.quit)
		return nil
	}
	log.Debug("handlers stop")
	return errors.New("already stopped")
}

// OnStop implements Service by doing nothing.
// NOTE: Do not put anything in here,
// that way users don't need to call BaseService.OnStop()
func (bs *BaseService) OnStop() {}

// Reset implements Service by calling OnReset callback (if defined). An error
// will be returned if the handlers is running.
func (bs *BaseService) Reset() error {
	if !atomic.CompareAndSwapUint32(&bs.stopped, 1, 0) {
		log.Debug("handlers reset")
		return fmt.Errorf("can't reset running %s", bs.name)
	}

	// whether we've started, we can reset
	atomic.CompareAndSwapUint32(&bs.started, 1, 0)

	bs.quit = make(chan struct{})
	return bs.impl.OnReset()
}

// OnReset implements Service by panicking.
func (bs *BaseService) OnReset() error {
	panic("The handlers cannot be reset")
}

// IsRunning implements Service by returning true or false depending on the
// handler's state.
func (bs *BaseService) IsRunning() bool {
	return atomic.LoadUint32(&bs.started) == 1 && atomic.LoadUint32(&bs.stopped) == 0
}

// Wait blocks until the handlers is stopped.
func (bs *BaseService) Wait() {
	<-bs.quit
}

// String implements Service by returning a string representation of the handlers.
func (bs *BaseService) String() string {
	return bs.name
}

// Quit Implements Service by returning a quit channel.
func (bs *BaseService) Quit() <-chan struct{} {
	return bs.quit
}
