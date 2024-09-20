package errortypes

// IsRetryableError is retry-opt-out. Meaning we rather retry than fail (possibly for a long time or forever).
func IsRetryableError(err error) bool {
	if IsBadUserInput(err) {
		return false
	}

	return true
}
