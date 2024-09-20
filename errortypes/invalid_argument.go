package errortypes

import (
	"fmt"

	"github.com/pkg/errors"
)

type InvalidArgument struct {
	what string
}

func NewInvalidArgument(what string) error {
	return errors.WithStack(InvalidArgument{what: what})
}

func NewInvalidArgumentf(format string, a ...interface{}) error {
	return errors.WithStack(InvalidArgument{what: fmt.Sprintf(format, a...)})
}

func (b InvalidArgument) Error() string {
	return b.what
}

func IsInvalidArgument(err error) bool {
	return errors.Is(err, InvalidArgument{})
}

func (b InvalidArgument) Is(err error) bool {
	_, ok := err.(InvalidArgument)
	return ok
}
