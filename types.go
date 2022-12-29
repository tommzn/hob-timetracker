package timetracker

import "time"

// RecordType defines with kind of event has been tracked.
type RecordType string

const (

	// WORKDAY is used to track start and end of a usual workday.
	WORKDAY RecordType = "workday"

	// ILLNESS is used to track sick leave.
	ILLNESS RecordType = "illness"

	// VACATION to track holiday absence.
	VACATION RecordType = "vacation"

	// WEEKEND used for non-working days in a week.
	WEEKEND RecordType = "weekend"
)

// MonthlyReport included total amount of work for a month and details about each single day.
type MonthlyReport struct {

	// Year this report belongs to.
	Year int

	// Month this reprt has been created for.
	Month int

	// Location a report should be generated for.
	Location Locale

	// Days is the list of days in a momth.
	Days []Day

	// TotalWorkingTine is the entire working time of a month.
	TotalWorkingTime time.Duration
}

// TimeTrackingReport os a single captured time tracking event.
type TimeTrackingRecord struct {

	// Key is an unique identifier of a time tracking record.
	Key string

	// DeviceId is an identifier of a device which captures a time tracking record.
	DeviceId string

	// Type of a time tracking event.
	Type RecordType

	// Timestamp is the point in time a time tracking event has occurred.
	Timestamp time.Time

	// Estimated time tracking report a used to fill missing events. e.g. workday end if it not has been captured.
	Estimated bool
}

// Date is a single calendar day.
type Date struct {

	// Year this date belongs to.
	Year int

	// Month this date belongs to.
	Month int

	// Day this date belongs to.
	Day int
}

// Day is a single day of working, illness or vacations.
// It contains all time tracking events occurred for this date and calculated working/break time based on this events.
type Day struct {

	// Date of this day.
	Date

	// Type of a time tracking event.
	Type RecordType

	// WorkingTime is the total time of work for a day.
	WorkingTime time.Duration

	// BreakTime is total time of breaks for a day.
	BreakTime time.Duration

	// Events is a list of captured time tracking events.
	Events []TimeTrackingRecord
}

// Locale contains settings like country, region, time zone and working breaks.
type Locale struct {

	// ISO 3166-1 country code.
	Country string

	// Timezone, used to format time in reports.
	Timezone *string

	// DateFormat, used tp write dates in given format to report outputs.
	DateFormat *string

	// DefaultWorkTime is used if there's no end of work for a day, including breaks.
	DefaultWorkTime time.Duration

	// Breaks is a map of working durations and breaks which have to be applied for this time.
	Breaks map[time.Duration]time.Duration
}

// Holiday is a single, public holiday.
type Holiday struct {

	// Date of this day.
	Date

	// Description of a public holiday.
	Description string
}
