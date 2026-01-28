package multierrgroup

import (
	"context"
	"sync"
)

// Collector is like Group but collects values of type T from the goroutines.
type Collector[T any] struct {
	vs []T
	m  sync.Mutex
	g  Group
}

// Go begins an operation in parallel that returns a value of type T.
func (r *Collector[T]) Go(f func() (T, error)) {
	r.g.Go(func() error {
		value, err := f()
		if err != nil {
			return err
		}
		r.m.Lock()
		defer r.m.Unlock()
		r.vs = append(r.vs, value)
		return nil
	})
}

// Wait waits for all routines to return and returns the collected values and
// any errors that occurred.
func (r *Collector[T]) Wait() ([]T, error) {
	err := r.g.Wait()
	return r.vs, err
}

// CollectorWithContext is like WithContext but returns a Collector to collect values
// of type T from the goroutines.
func CollectorWithContext[T any](ctx context.Context) (*Collector[T], context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Collector[T]{g: Group{cancel: cancel}}, ctx
}
