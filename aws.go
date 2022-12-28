package timetracker

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/ses"
)

// NewAWSConfig create a new config for AWS with passed region.
// In case nil is passed it try to access current region via AWS_REGION env variable.
func newAWSConfig(awsRegion *string) *aws.Config {

	if awsRegion != nil {
		return aws.NewConfig().WithRegion(*awsRegion)
	} else {
		if awsRegion, ok := os.LookupEnv("AWS_REGION"); ok {
			return aws.NewConfig().WithRegion(awsRegion)
		}
	}
	return aws.NewConfig()
}

// NewS3Downloader creates a new S3 download with passed region.
func newS3Downloader(awsRegion *string) *s3manager.Downloader {
	return s3manager.NewDownloader(newAwsSession(awsRegion))
}

// NewS3Uploader creates a new S3 uploader with passed region.
func newS3Uploader(awsRegion *string) *s3manager.Uploader {
	return s3manager.NewUploader(newAwsSession(awsRegion))
}

// NewS3Client create a new client to interact with aWS S3 buckets and objects.
func newS3Client(awsRegion *string) *s3.S3 {
	return s3.New(newAwsSession(awsRegion))
}

// NewAwsSession creates a new session for AWS connections.
func newAwsSession(awsRegion *string) *session.Session {
	return session.Must(session.NewSession(newAWSConfig(awsRegion)))
}

// newSESClient creates a new cloent to send emails via AWS Simple Mail Service.
func newSESClient(awsRegion *string) *ses.SES {
	return ses.New(newAwsSession(awsRegion))
}
