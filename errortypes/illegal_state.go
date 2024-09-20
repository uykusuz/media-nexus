package errortypes

import (
	"fmt"

	"github.com/pkg/errors"
)

type IllegalState struct {
	what string
}

func NewIllegalStatef(format string, a ...interface{}) error {
	return errors.WithStack(IllegalState{what: fmt.Sprintf(format, a...)})
}

func (b IllegalState) Error() string {
	return b.what
}

func IsIllegalState(err error) bool {
	return errors.Is(err, IllegalState{})
}

func (b IllegalState) Is(err error) bool {
	_, ok := err.(IllegalState)
	return ok
}
