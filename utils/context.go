package utils

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type SignalContext struct {
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	activeTasks sync.Map
}

// Return cancelable context
func (c *SignalContext) Context() context.Context {
	return c.ctx
}

func (c *SignalContext) Wait() {
	c.wg.Wait()
}

func (c *SignalContext) Add(v int) {
	c.wg.Add(v)
}

func (c *SignalContext) Done() {
	c.wg.Done()
}

// Manage Active Tasks
func (c *SignalContext) AddTask(id string, canceler context.CancelFunc) {
	c.activeTasks.Store(id, canceler)
	c.Add(1)
}

func (c *SignalContext) CancelTask(id string) {
	if cancel, exists := c.activeTasks.Load(id); exists {
		cancel.(context.CancelFunc)()
		c.activeTasks.Delete(id)
		c.Done()
	}
}

func (c *SignalContext) IsActiveTask(id string) bool {
	_, exists := c.activeTasks.Load(id)
	return exists
}

func NewSignalContext() *SignalContext {
	ctx, cancel := context.WithCancel(context.Background())

	context := &SignalContext{
		ctx:    ctx,
		cancel: cancel,
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		context.cancel()
	}()

	return context
}
