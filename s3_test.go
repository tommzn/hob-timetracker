package timetracker

import (
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
	"time"
)

type S3TestSuite struct {
	suite.Suite
}

func TestS3TestSuite(t *testing.T) {
	suite.Run(t, new(S3TestSuite))
}

func (suite *S3TestSuite) skipCI() {
	if _, isSet := os.LookupEnv("CI"); isSet {
		suite.T().Skip("Skip test in CI environment.")
	}
}

func (suite *S3TestSuite) TestCapture() {

	suite.skipCI()

	repo := suite.s3RepoForTest()
	deviceId := deviceIdForTest()

	suite.Nil(repo.Capture(deviceId, WORKDAY))
	suite.Nil(repo.Captured(deviceId, WORKDAY, time.Now().Add(1*time.Hour)))
}

func (suite *S3TestSuite) TestListRecords() {

	suite.skipCI()

	repo := suite.s3RepoForTest()
	deviceId := deviceIdForTest()

	suite.Nil(repo.Capture(deviceId, WORKDAY))
	start := time.Now()
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	end := start.AddDate(0, 0, 1).Add(-1 * time.Second)
	records, err := repo.ListRecords(deviceId, start, end)
	suite.Nil(err)
	suite.True(len(records) >= 1)

	records1, err1 := repo.ListRecords(deviceId, time.Now().Add(1*time.Minute), time.Now())
	suite.NotNil(err1)
	suite.Len(records1, 0)
}

func (suite *S3TestSuite) TestRecordCrudActions() {

	suite.skipCI()

	repo := suite.s3RepoForTest()

	record := TimeTrackingRecord{
		DeviceId:  "Device01",
		Type:      WORKDAY,
		Timestamp: time.Now(),
	}

	record1, err := repo.Add(record)
	suite.Nil(err)
	suite.True(len(record1.Key) > 0)

	suite.Nil(repo.Delete(record1.Key))
}

func (suite *S3TestSuite) TestPublishReport() {

	suite.skipCI()

	publisher := suite.s3PublisherForTest()
	suite.Nil(publisher.Send([]byte("Test-Report"), time.Now().String()))
}

func (suite *S3TestSuite) s3RepoForTest() *S3Repository {
	bucket, ok := os.LookupEnv("AWS_S3_TEST_BUCKET")
	suite.True(ok)
	path := "timetracker-test"
	return NewS3Repository(nil, &bucket, &path)
}

func (suite *S3TestSuite) s3PublisherForTest() *S3Publisher {
	bucket, ok := os.LookupEnv("AWS_S3_TEST_BUCKET")
	suite.True(ok)
	path := "timetracker-reports-test"
	return NewS3Publisher(nil, &bucket, &path, loggerForTest())
}
