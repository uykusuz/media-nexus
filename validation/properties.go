package validation

import (
	"math"
	"net/url"

	"media-nexus/errortypes"
)

func IsValidIntProperty(container string, name string, value int, min int, max int) error {
	if value < min {
		return errortypes.NewBadUserInputf("%v.%v less than %d", container, name, min)
	}

	if value > max {
		return errortypes.NewBadUserInputf("%v.%v greater than %d", container, name, max)
	}

	return nil
}

func IsValidPortProperty(container string, name string, value int) error {
	return IsValidIntProperty(container, name, value, 1, math.MaxInt16)
}

func IsValidURLProperty(container string, name string, value string) (*url.URL, error) {
	url, err := url.Parse(value)
	if err != nil {
		return nil, errortypes.NewBadUserInputf("%v.%v: '%v' is not a valid url: %v", container, name, value, err)
	}

	return url, nil
}

func IsValidStringProperty(containerType string, propertyName string, value string) error {
	if len(value) < 1 {
		return errortypes.NewBadUserInputf("%v in %v empty, but should not", propertyName, containerType)
	}

	return nil
}
