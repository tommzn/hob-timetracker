package timetracker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type LocalRepositoryTestSuite struct {
	suite.Suite
}

func TestLocalRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(LocalRepositoryTestSuite))
}

func (suite *LocalRepositoryTestSuite) TestCapture() {

	repo := NewLocaLRepository()
	t1 := time.Now()
	d1 := asDate(t1)
	deviceId := deviceIdForTest()

	repo.Capture(deviceId, WORKDAY)
	suite.Len(repo.Records, 1)
	r1, ok1 := repo.Records[deviceId][d1]
	suite.True(ok1)
	suite.Len(r1, 1)
	suite.Equal(WORKDAY, r1[0].Type)

	repo.Capture(deviceId, WORKDAY)
	suite.Len(repo.Records, 1)
	r1_1, ok1_1 := repo.Records[deviceId][d1]
	suite.True(ok1_1)
	suite.Len(r1_1, 2)
	suite.Equal(WORKDAY, r1_1[1].Type)

	t1 = t1.Add(123 * 24 * time.Hour)
	repo.Captured(deviceId, ILLNESS, t1)
	suite.Len(repo.Records, 1)
	suite.Len(repo.Records[deviceId], 2)
	d2 := asDate(t1)
	r2, ok2 := repo.Records[deviceId][d2]
	suite.True(ok2)
	suite.Len(r2, 1)
	suite.Equal(ILLNESS, r2[0].Type)
}

func (suite *LocalRepositoryTestSuite) TestListRecords() {

	deviceId := deviceIdForTest()
	repo := NewLocaLRepository()
	prepareRecords(repo, deviceId)

	records, err := repo.ListRecords(deviceId, time.Now(), time.Now().Add(3*24*time.Hour))
	suite.Nil(err)
	suite.Len(records, 12)

	records2, err2 := repo.ListRecords(deviceId, time.Now().Add(5*24*time.Hour), time.Now().Add(5*24*time.Hour))
	suite.Nil(err2)
	suite.Len(records2, 1)

	records3, err3 := repo.ListRecords(deviceId, time.Now().Add(2*24*time.Hour), time.Now())
	suite.NotNil(err3)
	suite.Len(records3, 0)
}

func (suite *LocalRepositoryTestSuite) TestRecordCrudActions() {

	repo := NewLocaLRepository()
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

func prepareRecords(repo *LocaLRepository, deviceId string) {

	durations := []time.Duration{
		time.Duration(0),
		1 * time.Minute,
		60 * time.Minute,
		24*time.Hour + time.Duration(0),
		24*time.Hour + 1*time.Minute,
		24*time.Hour + 60*time.Minute,
		2*24*time.Hour + time.Duration(0),
		2*24*time.Hour + 1*time.Minute,
		2*24*time.Hour + 60*time.Minute,
		3*24*time.Hour + time.Duration(0),
		3*24*time.Hour + 1*time.Minute,
		3*24*time.Hour + 60*time.Minute,
		4 * 24 * time.Hour,
		5 * 24 * time.Hour,
		6 * 24 * time.Hour,
	}
	t := time.Now()
	for _, duration := range durations {
		repo.Captured(deviceId, WORKDAY, t.Add(duration))
	}
}
