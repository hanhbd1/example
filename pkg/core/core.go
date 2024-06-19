package core

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"example/pkg/log"
)

const (
	defaultGracefulTimeout = time.Second * 5
	defaultMaxProcs        = 300
)

type App struct {
	name            string
	gracefulTimeout time.Duration
	baseCtx         context.Context
	maxProcs        int

	startFunc    func(ctx context.Context)
	shutdownFunc func(ctx context.Context)
}

type AppAction interface {
	Start(ctx context.Context)
	Shutdown(ctx context.Context)
}

// NewApp init app instance with options
func NewApp(options ...AppOption) *App {
	a := &App{
		baseCtx: context.Background(),
	}

	for _, option := range options {
		option(a)
	}

	if a.gracefulTimeout == 0 {
		a.gracefulTimeout = defaultGracefulTimeout
	}
	if a.maxProcs <= 0 {
		a.maxProcs = defaultMaxProcs
	}
	// set max goroutines
	runtime.GOMAXPROCS(a.maxProcs)
	return a
}

func (a *App) Run() {
	// Interrupt handler
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		select {
		case s := <-c:
			log.Infof("received signal: %s", s)
		case <-a.baseCtx.Done():
			log.Info("received context done")
		}

		// Create a deadline to wait for.
		closeCtx, cancel := context.WithTimeout(context.Background(), a.gracefulTimeout)
		defer cancel()
		a.shutdownFunc(closeCtx)
	}()

	a.startFunc(a.baseCtx)
}
