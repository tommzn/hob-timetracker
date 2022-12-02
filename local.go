package timetracker

import (
	"fmt"
	"time"
)

// NewLocaLRepository create a new, in menory, time tracker.
func NewLocaLRepository() *LocaLRepository {
	return &LocaLRepository{Records: make(map[string]map[Date][]TimeTrackingRecord)}
}

// LocaLRepository is an in memory time tracker.
type LocaLRepository struct {
	Records map[string]map[Date][]TimeTrackingRecord
}

// Capture will create a time tracking record with passed type at time this method has been called.
func (repo *LocaLRepository) Capture(deviceId string, recordType RecordType) error {
	return repo.Captured(deviceId, recordType, time.Now())
}

// Captured creates a time tracking record for passed point in time.
func (repo *LocaLRepository) Captured(deviceId string, recordType RecordType, timestamp time.Time) error {
	timestamp = timestamp.UTC()
	date := asDate(timestamp)
	if _, ok := repo.Records[deviceId]; !ok {
		repo.Records[deviceId] = make(map[Date][]TimeTrackingRecord)
	}
	if _, ok := repo.Records[deviceId][date]; !ok {
		repo.Records[deviceId][date] = []TimeTrackingRecord{}
	}
	repo.Records[deviceId][date] = append(repo.Records[deviceId][date], TimeTrackingRecord{DeviceId: deviceId, Type: recordType, Timestamp: timestamp, Estimated: false})
	return nil
}

// ListRecords returns available time tracking records for given range.
func (repo *LocaLRepository) ListRecords(deviceId string, start time.Time, end time.Time) ([]TimeTrackingRecord, error) {

	records := []TimeTrackingRecord{}
	if end.Before(start) {
		return records, fmt.Errorf("Invalid range: %s - %s", start, end)
	}
	if deviceREcords, ok := repo.Records[deviceId]; ok {
		for isDayBeforeOrEqual(start, end) {
			if recordsForDay, ok := deviceREcords[asDate(start)]; ok {
				records = append(records, recordsForDay...)
			}
			start = nextDay(start)
		}
	}
	return records, nil
}
