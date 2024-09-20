package errortypes

import (
	"fmt"

	"github.com/pkg/errors"
)

type ResourceAlreadyExists struct {
	what string
}

func NewResourceAlreadyExists(resourceID string) error {
	return NewResourceAlreadyExistsf("resource with id '%s' already exists", resourceID)
}

func NewResourceAlreadyExistsWithMessage(message string) error {
	return errors.WithStack(ResourceAlreadyExists{what: message})
}

func NewResourceAlreadyExistsf(format string, a ...interface{}) error {
	what := fmt.Sprintf(format, a...)
	return errors.WithStack(ResourceAlreadyExists{what: what})
}

func (b ResourceAlreadyExists) Error() string {
	return b.what
}

func IsResourceAlreadyExists(err error) bool {
	return errors.Is(err, ResourceAlreadyExists{})
}

func (b ResourceAlreadyExists) Is(err error) bool {
	_, ok := err.(ResourceAlreadyExists)
	return ok
}
