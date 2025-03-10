package pkg

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"sync"
)

type (
	ErrGroupWithRecovery struct {
		parent    *errgroup.Group
		err       error
		ctxCancel context.CancelFunc
		ctx       context.Context
		mt        *sync.Mutex
	}

	ErrGroupWithRecoveryAndSharedMutex struct {
		parent      *ErrGroupWithRecovery
		sharedMutex *sync.Mutex
	}
)

func NewErrGroupWithRecoveryAndSharedMutex(ctx context.Context) *ErrGroupWithRecoveryAndSharedMutex {
	return &ErrGroupWithRecoveryAndSharedMutex{
		parent:      NewErrGroupWithRecovery(ctx),
		sharedMutex: &sync.Mutex{},
	}
}

func NewErrGroupWithRecovery(ctx context.Context) *ErrGroupWithRecovery {
	ctx, cancel := context.WithCancel(ctx)
	parent, ctx := errgroup.WithContext(ctx)
	return &ErrGroupWithRecovery{
		parent:    parent,
		ctxCancel: cancel,
		ctx:       ctx,
		mt:        &sync.Mutex{},
	}
}

func (g *ErrGroupWithRecovery) Go(f func() error) {
	g.parent.Go(func() error {
		defer func() {
			if p := recover(); p != nil {
				g.setErr(fmt.Errorf("panic: %v", p))
				log.Printf("ErrGroupWithRecovery panic")
				g.ctxCancel()
			}
		}()

		select {
		case <-g.ctx.Done():
			return g.ctx.Err()
		default:
			if err := f(); err != nil {
				g.setErr(err)
				g.ctxCancel()
			}
		}

		return g.err
	})
}

func (g *ErrGroupWithRecovery) setErr(err error) {
	if err != nil {
		g.mt.Lock()
		g.err = err
		g.mt.Unlock()
	}
}

func (g *ErrGroupWithRecovery) Wait() error {
	go func() {
		if err := g.parent.Wait(); err != nil {
			g.setErr(err)
		}
	}()

	<-g.ctx.Done()

	if g.ctx.Err() != context.Canceled {
		return g.ctx.Err()
	}

	return g.err
}

func (g *ErrGroupWithRecoveryAndSharedMutex) Go(f func(m *sync.Mutex) error) {
	g.parent.Go(func() error {
		return f(g.sharedMutex)
	})
}

func (g *ErrGroupWithRecoveryAndSharedMutex) Wait() error {
	return g.parent.Wait()
}

func GoWithRecovery(f func()) {
	go func() {
		defer func() {
			if p := recover(); p != nil {
				log.Printf("GoWithRecovery: panic: %s", fmt.Sprintf("%v", p))
			}
		}()

		f()
	}()
}
