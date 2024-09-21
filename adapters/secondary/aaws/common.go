package aaws

import "github.com/aws/aws-sdk-go-v2/service/s3"

func WithRegion(region string) func(o *s3.Options) {
	return func(o *s3.Options) {
		o.Region = region
	}
}
