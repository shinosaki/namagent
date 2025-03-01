package signalcontext

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type SignalContext struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func NewSignalContext() *SignalContext {
	ctx, cancel := context.WithCancel(context.Background())

	sc := &SignalContext{
		ctx:    ctx,
		cancel: cancel,
	}

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		sc.Cancel()
	}()

	return sc
}

func (c *SignalContext) Context() context.Context {
	return c.ctx
}

func (c *SignalContext) Cancel() {
	c.cancel()
}

func (c *SignalContext) SelfContext() context.Context {
	return context.WithValue(c.ctx, SELF, c)
}

func (c *SignalContext) WithValue(key ContextKey, value any) *SignalContext {
	return &SignalContext{
		ctx:    context.WithValue(c.ctx, key, value),
		cancel: c.cancel,
	}
}

func (c *SignalContext) GetValue(key ContextKey) any {
	return c.ctx.Value(key)
}
