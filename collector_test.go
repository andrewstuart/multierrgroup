package multierrgroup_test

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"testing"

	"github.com/andrewstuart/multierrgroup"
)

func TestCollector_Go_Success(t *testing.T) {
	var r multierrgroup.Collector[int]

	r.Go(func() (int, error) { return 1, nil })
	r.Go(func() (int, error) { return 2, nil })
	r.Go(func() (int, error) { return 3, nil })

	values, err := r.Wait()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(values) != 3 {
		t.Fatalf("expected 3 values, got %d", len(values))
	}

	sort.Ints(values)
	if values[0] != 1 || values[1] != 2 || values[2] != 3 {
		t.Fatalf("unexpected values: %v", values)
	}
}

func TestCollector_Go_AllErrors(t *testing.T) {
	var r multierrgroup.Collector[int]

	r.Go(func() (int, error) { return 0, errors.New("error1") })
	r.Go(func() (int, error) { return 0, errors.New("error2") })

	values, err := r.Wait()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if len(values) != 0 {
		t.Fatalf("expected 0 values, got %d", len(values))
	}

	unw, ok := err.(interface{ Unwrap() []error })
	if !ok {
		t.Fatal("expected error to implement Unwrap() []error")
	}
	if len(unw.Unwrap()) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(unw.Unwrap()))
	}
}

func TestCollector_Go_PartialErrors(t *testing.T) {
	var r multierrgroup.Collector[int]

	r.Go(func() (int, error) { return 1, nil })
	r.Go(func() (int, error) { return 0, errors.New("error1") })
	r.Go(func() (int, error) { return 2, nil })

	values, err := r.Wait()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if len(values) != 2 {
		t.Fatalf("expected 2 values, got %d", len(values))
	}

	sort.Ints(values)
	if values[0] != 1 || values[1] != 2 {
		t.Fatalf("unexpected values: %v", values)
	}
}

func TestCollectorWithContext(t *testing.T) {
	ctx := context.Background()
	r, ctx := multierrgroup.CollectorWithContext[string](ctx)

	r.Go(func() (string, error) { return "a", nil })
	r.Go(func() (string, error) { return "b", nil })

	values, err := r.Wait()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(values) != 2 {
		t.Fatalf("expected 2 values, got %d", len(values))
	}

	sort.Strings(values)
	if values[0] != "a" || values[1] != "b" {
		t.Fatalf("unexpected values: %v", values)
	}
}

func TestCollectorWithContext_Cancellation(t *testing.T) {
	ctx := context.Background()
	r, ctx := multierrgroup.CollectorWithContext[int](ctx)

	r.Go(func() (int, error) { return 0, errors.New("fail") })

	_, err := r.Wait()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	select {
	case <-ctx.Done():
		// Context should be cancelled after error
	default:
		t.Fatal("expected context to be cancelled")
	}
}

func ExampleCollector() {
	var r multierrgroup.Collector[int]

	r.Go(func() (int, error) { return 1, nil })
	r.Go(func() (int, error) { return 2, nil })
	r.Go(func() (int, error) { return 3, nil })

	values, err := r.Wait()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	sort.Ints(values)
	fmt.Println(values)
	// Output: [1 2 3]
}
