package multierrgroup

import (
	"context"
	"errors"
	"sync"
)

// A Group can perform many operations in parallel, returning all errors that
// occur, as a hashicorp/go-multierror.Error instance,  once all operations
// have completed.
type Group struct {
	wg     sync.WaitGroup
	mut    sync.Mutex
	errs   []error
	cancel context.CancelFunc
}

// Go begins an operation in parallel.
func (g *Group) Go(f func() error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		if err := f(); err != nil {
			g.mut.Lock()
			g.errs = append(g.errs, err)
			if g.cancel != nil {
				g.cancel()
			}
			g.mut.Unlock()
		}
	}()
}

// GoWithContext is a helper for funcs or methods that already take a context,
// to avoid annoying boilerplate.
func (g *Group) GoWithContext(ctx context.Context, f func(context.Context) error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		if err := f(ctx); err != nil {
			g.mut.Lock()
			g.errs = append(g.errs, err)
			if g.cancel != nil {
				g.cancel()
			}
			g.mut.Unlock()
		}
	}()
}

// Wait waits for all routines to return and returns the errors, if any.
func (g *Group) Wait() error {
	g.wg.Wait()
	return errors.Join(g.errs...)
}

func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Group{cancel: cancel}, ctx
}
