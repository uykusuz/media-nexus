package errortypes

import (
	"fmt"

	"github.com/pkg/errors"
)

// UpstreamCommunicationError is an error that results in a 502 response.
type UpstreamCommunicationError struct {
	upstreamName string
	reason       string
}

func NewUpstreamCommunicationError(upstreamName string, reason error) error {
	return errors.WithStack(UpstreamCommunicationError{
		upstreamName: upstreamName,
		reason:       fmt.Sprintf("%v", reason),
	})
}

func NewUpstreamCommunicationErrorf(upstreamName string, format string, a ...interface{}) error {
	return errors.WithStack(UpstreamCommunicationError{
		upstreamName: upstreamName,
		reason:       fmt.Sprintf(format, a...),
	})
}

func (e UpstreamCommunicationError) Error() string {
	return fmt.Sprintf("failed to communicate with endpoint %v: %v", e.upstreamName, e.reason)
}

func IsUpstreamCommunicationError(err error) bool {
	return errors.Is(err, UpstreamCommunicationError{})
}

func (e UpstreamCommunicationError) Is(err error) bool {
	_, ok := err.(UpstreamCommunicationError)
	return ok
}
