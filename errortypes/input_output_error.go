package errortypes

import (
	"fmt"

	"github.com/pkg/errors"
)

type InputOutputError struct {
	what string
}

func NewInputOutputError(what string) error {
	return errors.WithStack(InputOutputError{what: what})
}

func NewInputOutputErrorf(format string, a ...interface{}) error {
	return errors.WithStack(InputOutputError{what: fmt.Sprintf(format, a...)})
}

func (b InputOutputError) Error() string {
	return b.what
}

func IsInputOutputError(err error) bool {
	return errors.Is(err, InputOutputError{})
}

func (b InputOutputError) Is(err error) bool {
	_, ok := err.(InputOutputError)
	return ok
}
