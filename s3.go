package timetracker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	log "github.com/tommzn/go-log"
	utils "github.com/tommzn/go-utils"
)

// NewS3Publisher returns a new publisher to upload reports to AWS S3.
func NewS3Publisher(awsRegion, bucket, basePath *string, logger log.Logger) *S3Publisher {
	return &S3Publisher{
		bucket:   bucket,
		basePath: basePath,
		s3:       newS3Client(awsRegion),
		uploader: newS3Uploader(awsRegion),
		logger:   logger,
	}
}

// NewS3Repository create a new repository to store time tracking records in AWS S3.
func NewS3Repository(awsRegion, bucket, basePath *string) *S3Repository {
	return &S3Repository{
		bucket:     bucket,
		basePath:   basePath,
		s3:         newS3Client(awsRegion),
		downloader: newS3Downloader(awsRegion),
		uploader:   newS3Uploader(awsRegion),
	}
}

// S3Publisher uploads given report to an AWS S3 bucket.
type S3Publisher struct {
	bucket   *string
	basePath *string
	s3       *s3.S3
	uploader *s3manager.Uploader
	logger   log.Logger
}

// S3Repository uses AWS S3 bucket to persist time tracking records.
type S3Repository struct {
	bucket     *string
	basePath   *string
	s3         *s3.S3
	downloader *s3manager.Downloader
	uploader   *s3manager.Uploader
}

// Capture will create a time tracking record with passed type at time this method has been called.
func (repo *S3Repository) Capture(deviceId string, recordType RecordType) error {
	return repo.Captured(deviceId, recordType, time.Now())
}

// Captured creates a time tracking record for passed point in time.
func (repo *S3Repository) Captured(deviceId string, recordType RecordType, timestamp time.Time) error {

	timeTrackingRecord := TimeTrackingRecord{
		DeviceId:  deviceId,
		Type:      recordType,
		Timestamp: timestamp.UTC().Round(time.Second),
	}
	objectPath := repo.newS3ObjectPath(deviceId, timeTrackingRecord.Timestamp)
	timeTrackingRecord.Key = *objectPath + repo.newRecordId()
	content, _ := json.Marshal(timeTrackingRecord)
	uploadInput := &s3manager.UploadInput{
		Bucket: repo.bucket,
		Key:    aws.String(timeTrackingRecord.Key),
		Body:   bytes.NewReader(content),
	}
	_, uploadErr := repo.uploader.Upload(uploadInput)
	return uploadErr
}

// ListRecords returns all records captured for given device id and time range.
func (repo *S3Repository) ListRecords(deviceId string, start time.Time, end time.Time) ([]TimeTrackingRecord, error) {

	records := []TimeTrackingRecord{}
	if start.After(end) {
		return records, fmt.Errorf("Invalid range: %s - %s", start, end)
	}

	start = start.UTC().Round(time.Second)
	end = end.UTC().Round(time.Second)

	objectKeys := []*string{}
	listObjectsInput := &s3.ListObjectsInput{
		Bucket: repo.bucket,
	}

	currentDate := start
	for isDayBeforeOrEqual(currentDate, end) {

		listObjectsInput.Prefix = repo.newS3ObjectPath(deviceId, currentDate)
		listObjectsOutput, err := repo.s3.ListObjects(listObjectsInput)
		if err != nil {
			return records, err
		}
		for _, s3object := range listObjectsOutput.Contents {
			objectKeys = append(objectKeys, s3object.Key)
		}
		currentDate = nextDay(currentDate)
	}

	for _, key := range objectKeys {

		requestInput := &s3.GetObjectInput{
			Bucket: repo.bucket,
			Key:    key,
		}

		buf := new(aws.WriteAtBuffer)
		_, err := repo.downloader.Download(buf, requestInput)
		if err != nil {
			return records, err
		}

		timeTrackingRecord := &TimeTrackingRecord{}
		decodeErr := json.NewDecoder(bytes.NewReader(buf.Bytes())).Decode(timeTrackingRecord)
		if decodeErr != nil {
			return records, decodeErr
		}
		timeTrackingRecord.Key = *key
		if isInRange(start, end, timeTrackingRecord.Timestamp) {
			records = append(records, *timeTrackingRecord)
		}

	}
	return records, nil
}

// Add creates a new time tracking record with given values. Same time tacking record will be
// returned together with a generated key.
func (repo *S3Repository) Add(record TimeTrackingRecord) (TimeTrackingRecord, error) {

	record.Timestamp = record.Timestamp.UTC().Round(time.Second)
	objectPath := repo.newS3ObjectPath(record.DeviceId, record.Timestamp)
	record.Key = *objectPath + repo.newRecordId()

	content, _ := json.Marshal(record)
	uploadInput := &s3manager.UploadInput{
		Bucket: repo.bucket,
		Key:    aws.String(record.Key),
		Body:   bytes.NewReader(content),
	}
	_, err := repo.uploader.Upload(uploadInput)
	return record, err
}

// Delete will remove given time tracking record.
func (repo *S3Repository) Delete(key string) error {

	deleteObjectInput := &s3.DeleteObjectInput{
		Bucket: repo.bucket,
		Key:    aws.String(key),
	}

	_, err := repo.s3.DeleteObject(deleteObjectInput)
	return err
}

// NewS3ObjectPath create a S3 object key including passed device id and date.
// Will add a path prefix if it has been defined at creating this repository.
func (repo *S3Repository) newS3ObjectPath(deviceId string, t time.Time) *string {
	utc := t.UTC()
	path := fmt.Sprintf("/%s/%04d/%02d/%02d/", deviceId, utc.Year(), utc.Month(), utc.Day())
	if repo.basePath != nil {
		path = *repo.basePath + path
	}
	return &path
}

// NewRecordId generates a new UUID v4.
func (repo *S3Repository) newRecordId() string {
	return utils.NewId()
}

// Send will upload given report data to AWS S3.
func (publisher *S3Publisher) Send(data []byte, name string) error {
	uploadInput := &s3manager.UploadInput{
		Bucket: publisher.bucket,
		Key:    publisher.objectKey(name),
		Body:   bytes.NewReader(data),
	}
	if _, uploadErr := publisher.uploader.Upload(uploadInput); uploadErr != nil {
		publisher.logger.Error("Unable to upload file to S3, reason: ", uploadErr)
		return uploadErr
	}
	publisher.logger.Debugf("Successfully uploaded to S3 at %s/%s", *uploadInput.Bucket, *uploadInput.Key)
	return nil
}

// ObjectKey creates a S3 object key for given report name,
// Will add a path prefix if it has been defined at creating this publisher.
func (publisher *S3Publisher) objectKey(name string) *string {
	if publisher.basePath != nil {
		if !strings.HasSuffix(*publisher.basePath, "/") {
			name = "/" + name
		}
		name = *publisher.basePath + name
	}
	return &name
}
