package errortypes

import (
	"fmt"

	"github.com/pkg/errors"
)

type ResourceNotFound struct {
	what string
}

func NewResourceNotFound(resourceID string) error {
	what := fmt.Sprintf("resource '%s' not found", resourceID)
	return errors.WithStack(ResourceNotFound{what: what})
}

func NewResourceNotFoundf(format string, a ...interface{}) error {
	return NewResourceNotFound(fmt.Sprintf(format, a...))
}

func (b ResourceNotFound) Error() string {
	return b.what
}

func IsResourceNotFound(err error) bool {
	return errors.Is(err, ResourceNotFound{})
}

func (b ResourceNotFound) Is(err error) bool {
	_, ok := err.(ResourceNotFound)
	return ok
}

func IsOnlyResourceNotFound(errs []error) bool {
	for _, err := range errs {
		if !IsResourceNotFound(err) {
			return false
		}
	}

	return true
}
