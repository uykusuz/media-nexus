package amongodb

import (
	"media-nexus/errortypes"

	"go.mongodb.org/mongo-driver/mongo"
)

func handleError(err error) error {
	if err == nil {
		return nil
	}

	if mongo.IsDuplicateKeyError(err) {
		return errortypes.NewResourceAlreadyExistsf("resource already exists: %v", err)
	}

	if mongo.IsTimeout(err) {
		return errortypes.NewTimeoutf("timed out: %v", err)
	}

	if mongo.IsNetworkError(err) {
		return errortypes.NewServiceUnavailablef("network error: %v", err)
	}

	if err == mongo.ErrNoDocuments {
		return errortypes.NewResourceNotFound("mongo document not found")
	}

	return errortypes.NewUpstreamCommunicationErrorf("mongo", "failed to talk to mongodb: %v", err)
}
