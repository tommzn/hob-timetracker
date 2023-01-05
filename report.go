package timetracker

import (
	"sort"
	"time"
)

// NewReportCalulator returns a new calulator using given time tracking records and local.
func NewReportCalulator(records []TimeTrackingRecord, location Locale) *ReportCalulator {
	return &ReportCalulator{
		location: location,
		records:  records,
	}
}

// ReportCalulator calulates working and break time for all day given by time tracking records.
type ReportCalulator struct {

	// Location a report should be generated for.
	location Locale

	// Records os a list of all time tracking events a report should be generated for.
	records []TimeTrackingRecord
}

// WithTimeTrackingRecords will apply given records for calculation.
func (calculator *ReportCalulator) WithTimeTrackingRecords(records []TimeTrackingRecord) {
	calculator.records = records
}

// MonthlyReport generates a report for given year and month from existing time tracking records.
func (calculator *ReportCalulator) MonthlyReport(year, month int, latestType RecordType) (*MonthlyReport, error) {

	report := &MonthlyReport{
		Year:             year,
		Month:            month,
		Location:         calculator.location,
		Days:             []Day{},
		TotalWorkingTime: time.Duration(0),
	}

	days := splitToDays(calculator.records)
	for _, day := range days {
		if day.Date.Year == year && day.Date.Month == month {
			day.Type = calculator.determineTypeOf(day.Events)
			calculator.calculateWorkTimeForDay(&day)
			calculator.subtractBreaks(&day)
			report.TotalWorkingTime += day.WorkingTime
			report.Days = append(report.Days, day)
		}
	}
	fillVacationAndIllness(report, latestType)
	return report, nil
}

// GetEndOfWorkingDay will create an estimated time tracking record.
// If working time from passed records already exceeds given default working time end of working day
// will be one minute after last available timestamp.
// In other cases default working time will be added to first record to create end of working day.
func getEndOfWorkingDay(records []TimeTrackingRecord, defaultWorkTime time.Duration) TimeTrackingRecord {

	lastIdx := len(records) - 1
	totalWorkDuration := records[lastIdx].Timestamp.Sub(records[0].Timestamp)
	if totalWorkDuration < defaultWorkTime {
		return TimeTrackingRecord{Type: WORKDAY, Timestamp: records[0].Timestamp.Add(defaultWorkTime), Estimated: true}
	} else {
		return TimeTrackingRecord{Type: WORKDAY, Timestamp: records[lastIdx].Timestamp.Add(1 * time.Minute), Estimated: true}
	}
}

// DetermineTypeOf will analyze given records. In case an ILLNESS or VACATION records is present
// this type will be returned in all other case default value WORKDAY is returned.
func (calculator *ReportCalulator) determineTypeOf(records []TimeTrackingRecord) RecordType {

	if len(records) == 0 {
		return WORKDAY
	}

	for _, record := range records {
		if record.Type == ILLNESS || record.Type == VACATION {
			return record.Type
		}
	}
	return WORKDAY
}

// CalculateWorkTimeForDay summarizes total working time of goven day.
// In case an odd number of time tracking records is given it will add an extimated end of the day at first.
// Then it calculates duration between start and end pair of time tracking records and sum them up to
// a total working time for the day.
func (calculator *ReportCalulator) calculateWorkTimeForDay(day *Day) {

	if len(day.Events) == 0 {
		day.WorkingTime = 0
		day.BreakTime = 0
		return
	}

	sort.Slice(day.Events, func(i, j int) bool { return day.Events[i].Timestamp.Before(day.Events[j].Timestamp) })
	if len(day.Events)%2 != 0 {
		day.Events = append(day.Events, getEndOfWorkingDay(day.Events, calculator.location.DefaultWorkTime))
	}

	events := splitTimeTrackingRecords(day.Events, 2)
	for _, chunkOfEvents := range events {
		day.WorkingTime += chunkOfEvents[1].Timestamp.Sub(chunkOfEvents[0].Timestamp)
	}
}

// SubtractBreaks will reduce working time by given default breaks if there no breaks in working time
// or existing breaks doesn't reach default settings.
func (calculator *ReportCalulator) subtractBreaks(day *Day) {

	definedBreakTime := time.Duration(0)
	for workingTime, timeOfBreak := range calculator.location.Breaks {
		if day.WorkingTime >= workingTime {
			definedBreakTime += timeOfBreak
		}
	}

	day.BreakTime = day.Events[len(day.Events)-1].Timestamp.Sub(day.Events[0].Timestamp) - day.WorkingTime
	if day.BreakTime < definedBreakTime {
		day.WorkingTime -= definedBreakTime - day.BreakTime
		day.BreakTime = definedBreakTime
	}
}

// SplitToDays will walk trough given time tracking records and assign them to day of a month.
func splitToDays(records []TimeTrackingRecord) []Day {

	daysMap := make(map[Date]Day)
	for _, record := range records {

		date := asDate(record.Timestamp)
		day, ok := daysMap[date]
		if !ok {
			day = Day{
				Date:        date,
				WorkingTime: time.Duration(0),
				BreakTime:   time.Duration(0),
				Events:      []TimeTrackingRecord{},
			}
		}
		day.Events = append(day.Events, record)
		daysMap[date] = day
	}

	days := []Day{}
	for _, day := range daysMap {
		days = append(days, day)
	}
	return days
}

// SplitTimeTrackingRecords cuts given list of records into chunks with given size.
func splitTimeTrackingRecords(records []TimeTrackingRecord, chunkSize int) (chunks [][]TimeTrackingRecord) {
	for chunkSize < len(records) {
		records, chunks = records[chunkSize:], append(chunks, records[0:chunkSize:chunkSize])
	}
	return append(chunks, records)
}

// FillVacationAndIllness will add vacation and illness days.
// Thhis applies if a day with type ILLNESS or VACATION is available and following days doesn't exist in the list of days.
// For such cases days with same type will be generated until next day in the list or until the end of the month.
func fillVacationAndIllness(report *MonthlyReport, latestType RecordType) {

	sort.Slice(report.Days, func(i, j int) bool { return report.Days[i].Date.Before(report.Days[j].Date) })
	dayOfMonth := time.Date(report.Year, time.Month(report.Month), 1, 0, 0, 0, 0, time.UTC)

	days := []Day{}
	for _, day := range report.Days {

		if !day.Date.IsEqual(dayOfMonth) &&
			(latestType == ILLNESS || latestType == VACATION) {
			daysToFill := generateDays(dayOfMonth, day.Date, latestType)
			days = append(days, daysToFill...)
		} else {
			latestType = day.Type
		}
		dayOfMonth = time.Date(dayOfMonth.Year(), dayOfMonth.Month(), day.Date.Day, 0, 0, 0, 0, time.UTC)
		days = append(days, day)
		dayOfMonth = dayOfMonth.AddDate(0, 0, 1)
	}
	report.Days = fillToEndOfMonth(days)
}

// GenerateDays will create a list of empty days in given range and assign passed type to all of them.
func generateDays(startDay time.Time, endDay Date, recordType RecordType) []Day {
	days := []Day{}
	for startDay.Day() < endDay.Day {
		days = append(days, Day{Date: asDate(startDay), Type: recordType, WorkingTime: 0, BreakTime: 0})
		startDay = startDay.AddDate(0, 0, 1)
	}
	return days
}

// FillToEndOfMonth loop from last day in given list until end of month and fill vacation or illness days if required.
func fillToEndOfMonth(days []Day) []Day {

	if len(days) == 0 {
		return days
	}

	lastDayInList := days[len(days)-1]
	lastDayOfMonth := lastDayOfMonth(lastDayInList.Date)
	if (lastDayInList.Type != ILLNESS && lastDayInList.Type != VACATION) ||
		lastDayInList.Equal(lastDayOfMonth) {
		return days
	}
	daysToFill := generateDays(lastDayInList.Date.AsTime(), lastDayOfMonth, lastDayInList.Type)
	daysToFill = append(daysToFill, Day{Date: lastDayOfMonth, Type: lastDayInList.Type, WorkingTime: 0, BreakTime: 0})
	return append(days, daysToFill...)
}
