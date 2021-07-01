package multierrgroup_test

import (
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

	meg.Go(func() error {
		return fmt.Errorf("err1")
	})

	meg.Go(func() error {
		return &someError{msg: "foo"}
	})

	meg.Go(func() error {
		return nil
	})

	err := meg.Wait()

	fmt.Print(err.(*multierror.Error).Len(), ", ")
	var e *someError
	if ok := errors.As(err, &e); ok {
		fmt.Println(e)
	}

	// Output: 2, someerror: foo
}
