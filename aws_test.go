package timetracker

import (
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type AwsTestSuite struct {
	suite.Suite
}

func TestAwsTestSuite(t *testing.T) {
	suite.Run(t, new(AwsTestSuite))
}

func (suite *AwsTestSuite) TestAwsConfig() {

	currentAwsRegion := os.Getenv("AWS_REGION")

	expectedAwsRegion := "eu-central-5"
	os.Setenv("AWS_REGION", expectedAwsRegion)
	awsConfig1 := newAWSConfig(nil)
	suite.NotNil(awsConfig1)
	suite.Equal(expectedAwsRegion, *awsConfig1.Region)

	awsRegion := "eu-central-7"
	awsConfig2 := newAWSConfig(&awsRegion)
	suite.NotNil(awsConfig2)
	suite.Equal(awsRegion, *awsConfig2.Region)

	os.Unsetenv("AWS_REGION")
	awsConfig3 := newAWSConfig(nil)
	suite.NotNil(awsConfig3)
	suite.Nil(awsConfig3.Region)

	os.Setenv("AWS_REGION", currentAwsRegion)
}

func (suite *AwsTestSuite) TestS3Downloader() {

	awsRegion := "eu-central-5"
	downloader := newS3Downloader(&awsRegion)
	suite.NotNil(downloader)
}

func (suite *AwsTestSuite) TestS3Uploader() {

	awsRegion := "eu-central-5"
	downloader := newS3Uploader(&awsRegion)
	suite.NotNil(downloader)
}

func (suite *AwsTestSuite) TestS3Client() {

	awsRegion := "eu-central-5"
	client := newS3Client(&awsRegion)
	suite.NotNil(client)
}
