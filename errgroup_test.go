package multierrgroup_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/andrewstuart/multierrgroup"
	"github.com/hashicorp/go-multierror"
)

type someError struct {
	msg string
}

func (s *someError) Error() string {
	return fmt.Sprintf("someerror: %s", string(s.msg))
}

func ExampleMultierror() {
	var meg multierrgroup.Group
	ctx := context.Background()

	meg.Go(func() error {
		return fmt.Errorf("err1")
	})

	meg.Go(func() error {
		return &someError{msg: "foo"}
	})

	meg.Go(func() error {
		return nil
	})

	someFunc := func(context.Context) error {
		return nil
	}

	meg.GoWithContext(ctx, someFunc)

	err := meg.Wait()

	fmt.Print(err.(*multierror.Error).Len(), ", ")
	var e *someError
	if ok := errors.As(err, &e); ok {
		fmt.Println(e)
	}

	// Output: 2, someerror: foo
}
