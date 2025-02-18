package utils

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type ContextKey string

type SignalContext struct {
	wg     *sync.WaitGroup
	tasks  *sync.Map
	ctx    context.Context
	cancel context.CancelFunc
}

func NewSignalContext() *SignalContext {
	ctx, cancel := context.WithCancel(context.Background())

	sc := &SignalContext{
		wg:     &sync.WaitGroup{},
		tasks:  &sync.Map{},
		ctx:    ctx,
		cancel: cancel,
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		sc.cancel()
	}()

	return sc
}

func (c *SignalContext) Context() context.Context {
	return c.ctx
}

func (c *SignalContext) Wait() {
	c.wg.Wait()
}

func (c *SignalContext) IsActiveTask(id any) bool {
	_, ok := c.tasks.Load(id)
	return ok
}

func (c *SignalContext) AddTask(id any, canceler context.CancelFunc) {
	c.tasks.Store(id, canceler)
	c.wg.Add(1)
}

func (c *SignalContext) CancelTask(id any) {
	if canceler, ok := c.tasks.LoadAndDelete(id); ok {
		canceler.(context.CancelFunc)()
		c.wg.Done()
	}
}

func (c *SignalContext) GetValue(key ContextKey) any {
	return c.ctx.Value(key)
}

func (c *SignalContext) WithValue(key ContextKey, value any) *SignalContext {
	return &SignalContext{
		ctx:    context.WithValue(c.ctx, key, value),
		wg:     c.wg,
		tasks:  c.tasks,
		cancel: c.cancel,
	}
}
