package errortypes

import (
	"fmt"

	"github.com/pkg/errors"
)

type Timeout struct {
	what string
}

func NewTimeout(waitedFor string) error {
	return NewTimeoutf("timeout while waiting for %v", waitedFor)
}

func NewTimeoutf(format string, a ...interface{}) error {
	return errors.WithStack(Timeout{what: fmt.Sprintf(format, a...)})
}

func (b Timeout) Error() string {
	return b.what
}

func IsTimeout(err error) bool {
	return errors.Is(err, Timeout{})
}

func (b Timeout) Is(err error) bool {
	_, ok := err.(Timeout)
	return ok
}
