package timetracker

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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
	timestamp = timestamp.UTC().Round(time.Second)
	date := asDate(timestamp)
	if _, ok := repo.Records[deviceId]; !ok {
		repo.Records[deviceId] = make(map[Date][]TimeTrackingRecord)
	}
	if _, ok := repo.Records[deviceId][date]; !ok {
		repo.Records[deviceId][date] = []TimeTrackingRecord{}
	}
	record := TimeTrackingRecord{
		DeviceId:  deviceId,
		Type:      recordType,
		Timestamp: timestamp,
		Estimated: false,
	}
	record.Key = repo.recordKey(deviceId, record.Timestamp, len(repo.Records[deviceId][date]))
	repo.Records[deviceId][date] = append(repo.Records[deviceId][date], record)
	return nil
}

// ListRecords returns available time tracking records for given range.
func (repo *LocaLRepository) ListRecords(deviceId string, start time.Time, end time.Time) ([]TimeTrackingRecord, error) {

	start = start.Round(time.Second)
	end = end.Round(time.Second)

	records := []TimeTrackingRecord{}
	if end.Before(start) {
		return records, fmt.Errorf("Invalid range: %s - %s", start, end)
	}

	if deviceRecords, ok := repo.Records[deviceId]; ok {

		currentDate := start
		for isDayBeforeOrEqual(currentDate, end) {
			if recordsForDay, ok := deviceRecords[asDate(currentDate)]; ok {
				for idx, record := range recordsForDay {
					record.Key = repo.recordKey(deviceId, currentDate, idx)
					if isInRange(start, end, record.Timestamp) {
						records = append(records, record)
					}
				}
			}
			currentDate = nextDay(currentDate)
		}
	}
	return records, nil
}

// Add creates a new time tracking record with given values. Same time tacking record will be
// returned together with a generated key.
func (repo *LocaLRepository) Add(record TimeTrackingRecord) (TimeTrackingRecord, error) {

	record.Timestamp = record.Timestamp.UTC().Round(time.Second)
	deviceId := record.DeviceId
	date := asDate(record.Timestamp)
	if _, ok := repo.Records[deviceId]; !ok {
		repo.Records[deviceId] = make(map[Date][]TimeTrackingRecord)
	}
	if _, ok := repo.Records[deviceId][date]; !ok {
		repo.Records[deviceId][date] = []TimeTrackingRecord{}
	}
	record.Key = repo.recordKey(deviceId, record.Timestamp, len(repo.Records[deviceId][date]))
	repo.Records[deviceId][date] = append(repo.Records[deviceId][date], record)
	return record, nil
}

// Delete will remove given time tracking record.
func (repo *LocaLRepository) Delete(key string) error {

	keyParts := strings.Split(key, "/")
	if len(keyParts) != 3 {
		return errors.New("Invalid key passed.")
	}

	deviceId := keyParts[0]
	_, ok := repo.Records[deviceId]
	if !ok {
		return errors.New("Invalid deviceid: " + deviceId)
	}

	dateStr := keyParts[1]
	timestamp, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return err
	}

	date := asDate(timestamp)
	_, ok1 := repo.Records[deviceId][date]
	if !ok1 {
		return errors.New("Invalid date: " + dateStr)
	}

	idx, err := strconv.Atoi(keyParts[2])
	if err != nil {
		return err
	}

	if len(repo.Records[deviceId][date])-1 < idx {
		return errors.New("Invalid index: " + keyParts[2])
	}

	repo.Records[deviceId][date] = removeRecordAtIndex(repo.Records[deviceId][date], idx)
	return nil
}

// RecordKey compose given part to a record key.
func (repo *LocaLRepository) recordKey(deviceId string, day time.Time, idx int) string {
	return fmt.Sprintf("%s/%s/%d", deviceId, asDate(day).String(), idx)
}

// RemoveRecordAtIndex removes time tracking records from given list at specific index.
func removeRecordAtIndex(records []TimeTrackingRecord, index int) []TimeTrackingRecord {
	ret := make([]TimeTrackingRecord, 0)
	ret = append(ret, records[:index]...)
	return append(ret, records[index+1:]...)
}
