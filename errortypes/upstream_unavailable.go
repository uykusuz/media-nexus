package errortypes

import (
	"fmt"

	"github.com/pkg/errors"
)

type UpstreamUnavailable struct {
	what string
}

func NewUpstreamUnavailable(context string, err error) error {
	return errors.WithStack(
		UpstreamUnavailable{what: fmt.Sprintf("upstream unavailable for %v: %v", context, err)},
	)
}

func NewUpstreamUnavailablef(format string, a ...interface{}) error {
	return errors.WithStack(UpstreamUnavailable{
		what: fmt.Sprintf(format, a...),
	})
}

func (b UpstreamUnavailable) Error() string {
	return b.what
}

func IsUpstreamUnavailable(err error) bool {
	return errors.Is(err, UpstreamUnavailable{})
}

func (b UpstreamUnavailable) Is(err error) bool {
	_, ok := err.(UpstreamUnavailable)
	return ok
}
