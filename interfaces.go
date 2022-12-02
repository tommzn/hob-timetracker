package timetracker

import (
	"bytes"
	"time"
)

// TimeTracker is used to persist event, e.g start/end of a workday, illness and vacations.
// All timestamps are collected in UTC.
type TimeTracker interface {

	// Capture will create a time tracking record with passed type at time this method has been called.
	Capture(string, RecordType) error

	// Captured creates a time tracking record for passed point in time.
	Captured(string, RecordType, time.Time) error

	// ListRecords returns available time tracking records for given range.
	ListRecords(string, time.Time, time.Time) ([]TimeTrackingRecord, error)
}

// ReportCalculator creates a time tracking summary based on captured records.
type ReportCalculator interface {

	// MonthlyReport calculates a report for given year and month.
	MonthlyReport(int, int, RecordType) (*MonthlyReport, error)
}

// ReportFormatter generates an output for passed reports.
type ReportFormatter interface {

	// WriteMonthlyReportToFile will generate a report outout an writes it to given file.
	WriteMonthlyReportToFile(*MonthlyReport, string) error

	// WriteMonthlyReportToBuffer returns a buffer for gemerated report output.
	WriteMonthlyReportToBuffer(*MonthlyReport) (*bytes.Buffer, error)
}

// ReportPublisher sends given report to a defined target.
type ReportPublisher interface {

	// Send publishes given report data to a target.
	Send([]byte, string) error
}