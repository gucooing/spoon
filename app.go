package spoon

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type AppInfo interface {
}

type App struct {
	opts   options
	ctx    context.Context
	cancel context.CancelFunc
}

func New(opts ...Option) *App {
	o := options{
		ctx:  context.Background(),
		sigs: []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
	}
	for _, opt := range opts {
		opt(&o)
	}
	ctx, cancel := context.WithCancel(o.ctx)
	return &App{
		opts:   o,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (a *App) Run() error {
	sctx := NewContext(a.ctx, a)
	eg, ctx := errgroup.WithContext(sctx) // 全局
	wg := sync.WaitGroup{}

	// 启动服务
	octx := NewContext(a.opts.ctx, a)
	for _, srv := range a.opts.servers {
		eg.Go(func() error {
			<-ctx.Done()
			stopCtx := context.WithoutCancel(octx)
			return srv.Stop(stopCtx)
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			return srv.Start(octx)
		})
	}
	wg.Wait()

	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opts.sigs...)
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return nil
		case <-c:
			return a.Stop()
		}
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func (a *App) Stop() error {
	//sctx := NewContext(a.ctx, a)
	if a.cancel != nil {
		a.cancel()
	}
	return nil
}

type appKey struct{}

func NewContext(ctx context.Context, s AppInfo) context.Context {
	return context.WithValue(ctx, appKey{}, s)
}

func FromContext(ctx context.Context) (s AppInfo, ok bool) {
	s, ok = ctx.Value(appKey{}).(AppInfo)
	return
}
