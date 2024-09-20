package aaws

import (
	"context"
	"errors"
	"media-nexus/errortypes"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

func ensureBucketExists(ctx context.Context, client *s3.Client, bucketName string) error {
	exists, err := bucketExists(ctx, client, bucketName)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	return createBucket(ctx, client, bucketName)
}

func bucketExists(ctx context.Context, client *s3.Client, bucketName string) (bool, error) {
	_, err := client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err == nil {
		return true, nil
	}

	var apiError smithy.APIError
	if errors.As(err, &apiError) {
		switch apiError.(type) {
		case *types.NotFound:
			return false, nil
		}
	}

	return false, errortypes.NewUpstreamCommunicationErrorf(
		bucketName,
		"either we don't have access to bucket %v or another error occurred: %v",
		bucketName,
		err,
	)
}

func createBucket(ctx context.Context, client *s3.Client, bucketName string) error {
	_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return errortypes.NewUpstreamCommunicationErrorf(bucketName, "couldn't create bucket %v: %v", bucketName, err)
	}

	return nil
}
