package errortypes

import (
	"fmt"

	"github.com/pkg/errors"
)

type ServiceUnavailable struct {
	what string
}

func NewServiceUnavailable(service string) error {
	what := fmt.Sprintf("service unavailable: %v", service)
	return errors.WithStack(ServiceUnavailable{what: what})
}

func NewServiceUnavailablef(format string, a ...interface{}) error {
	what := fmt.Sprintf(format, a...)
	return errors.WithStack(ServiceUnavailable{what: what})
}

func (b ServiceUnavailable) Error() string {
	return b.what
}
