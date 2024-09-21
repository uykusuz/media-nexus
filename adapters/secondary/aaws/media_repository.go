package aaws

import (
	"context"
	"io"
	"media-nexus/errortypes"
	"media-nexus/ports"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type mediaRepository struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	bucket        string
}

func NewMediaRepository(client *s3.Client, presignClient *s3.PresignClient, bucket string) ports.MediaRepository {
	return &mediaRepository{client, presignClient, bucket}
}

func (r *mediaRepository) CreateMedia(ctx context.Context, key string, file io.Reader) error {
	err := ensureBucketExists(ctx, r.client, r.bucket)
	if err != nil {
		return err
	}

	uploadInput := &s3.PutObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
		Body:   file,
	}

	_, err = r.client.PutObject(ctx, uploadInput)
	if err != nil {
		return errortypes.NewInputOutputErrorf("failed to upload file to S3: %v", err)
	}

	return nil
}

func (r *mediaRepository) GetMediaUrl(ctx context.Context, key string, lifetime time.Duration) (string, error) {
	err := ensureBucketExists(ctx, r.client, r.bucket)
	if err != nil {
		return "", err
	}

	request, err := r.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = lifetime
	})
	if err != nil {
		return "", errortypes.NewUpstreamCommunicationErrorf(
			"couldn't get a presigned request to get %v:%v: %v",
			r.bucket,
			key,
			err,
		)
	}

	return request.URL, err
}

func (r *mediaRepository) DeleteAll(ctx context.Context, keys []string) error {
	err := ensureBucketExists(ctx, r.client, r.bucket)
	if err != nil {
		return err
	}

	objectIds := make([]types.ObjectIdentifier, 0, len(keys))
	for _, key := range keys {
		objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(key)})
	}

	input := &s3.DeleteObjectsInput{
		Bucket: aws.String(r.bucket),
		Delete: &types.Delete{
			Objects: objectIds,
		},
	}

	_, err = r.client.DeleteObjects(ctx, input)
	if err != nil {
		return errortypes.NewInputOutputErrorf("failed to delete objects from bucket: %v", err)
	}

	return nil
}
