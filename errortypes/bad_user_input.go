package errortypes

import (
	"fmt"

	"github.com/pkg/errors"
)

type BadUserInput struct {
	what string
}

func NewBadUserInput(what string) error {
	return errors.WithStack(BadUserInput{what: what})
}

func NewBadUserInputf(format string, a ...interface{}) error {
	return errors.WithStack(BadUserInput{what: fmt.Sprintf(format, a...)})
}

func (b BadUserInput) Error() string {
	return b.what
}

func IsBadUserInput(err error) bool {
	return errors.Is(err, BadUserInput{})
}

func (b BadUserInput) Is(err error) bool {
	_, ok := err.(BadUserInput)
	return ok
}
