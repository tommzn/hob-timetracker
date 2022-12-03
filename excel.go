package timetracker

import (
	"bytes"
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

// NewExcelReportFormatter returns a new formatter to generate an excel file for a report.
func NewExcelReportFormatter() *ExcelReportFormatter {
	return &ExcelReportFormatter{
		dateFormat: "2006-01-02",
		timeFormat: "15:04",
	}
}

// ExcelReportFormatter will generate Excel files for reports
type ExcelReportFormatter struct {

	// Holidays is a list of public holidays.
	// In case a day matches a date from this list it will be formatted with a specific backgound color.
	// see holidayStyleId for format infos.
	holidays map[Date]Holiday

	// dateFormat defnines the format a day should be printed in the outut.
	dateFormat string

	// TimeFormat defines output format for timestamps of time tracking records.
	timeFormat string

	// Timezone is used to convert timestamps from time tracking records, captured in UTC, to local time.
	timezone *time.Location

	// HeadlineStyleId, style generated for header columns.
	// Bold, with bottom border
	headlineStyleId int

	// weekendStyleId, style to format weekend days.
	// Light gray background color
	weekendStyleId int

	// DaysBottomStyleId, style for last day in a month.
	// Bottom border
	daysBottomStyleId int

	// HolidayStyleId, style to format holidays in outout. Overwrites weekend format.
	// Light green background
	holidayStyleId int

	// IllnessStyleId, style for days of illness.
	// Orange background
	illnessStyleId int

	// VacationStyleId, style for days of vacations.
	// Green background
	vacationStyleId int
}

// WithHolidays will assign give list of holidays for output formatting.
func (formatter *ExcelReportFormatter) WithHolidays(holidays []Holiday) {
	formatter.holidays = asHolidayMap(holidays)
}

// WriteMonthlyReportToFile will generate a report outout an writes it to given file.
func (formatter *ExcelReportFormatter) WriteMonthlyReportToFile(report *MonthlyReport, filename string) error {

	xls, err := formatter.generateOutput(report)
	if err != nil {
		return err
	}
	return xls.SaveAs(filename)
}

// WriteMonthlyReportToBuffer returns a buffer for gemerated report output.
func (formatter *ExcelReportFormatter) WriteMonthlyReportToBuffer(report *MonthlyReport) (*bytes.Buffer, error) {
	xls, err := formatter.generateOutput(report)
	if err != nil {
		return nil, err
	}
	return xls.WriteToBuffer()
}

// GenerateOutput writes entire report content, including all styles, to an excel file.
func (formatter *ExcelReportFormatter) generateOutput(report *MonthlyReport) (*excelize.File, error) {

	formatter.determineDateFormat(report)
	formatter.determineTimezone(report)
	days := generateIndexMap(report.Days)

	sheetName := fmt.Sprintf("%04d-%02d", report.Year, report.Month)
	xls := newExcelFile(sheetName)
	if err := formatter.createStyles(xls); err != nil {
		return nil, err
	}

	writeHeader(xls, sheetName)
	if err := xls.SetCellStyle(sheetName, getCellId("A", 1), getCellId("F", 1), formatter.headlineStyleId); err != nil {
		return nil, err
	}

	calendarDay := time.Date(report.Year, time.Month(report.Month), 1, 0, 0, 0, 0, time.UTC)
	lastDayOfMonth := calendarDay.AddDate(0, 1, -1)

	row := 2
	for isDayBeforeOrEqual(calendarDay, lastDayOfMonth) {

		date := asDate(calendarDay)
		if day, ok := days[date]; ok {
			formatter.appendRowForDay(day, xls, sheetName, row)
		} else {
			formatter.appendRowForDay(emptyDay(date), xls, sheetName, row)
		}
		row++
		calendarDay = nextDay(calendarDay)
	}

	if err := xls.SetCellStyle(sheetName, getCellId("A", row-1), getCellId("F", row-1), formatter.daysBottomStyleId); err != nil {
		return nil, err
	}
	writeSummary(xls, sheetName, row, report.TotalWorkingTime)
	xls.SetColWidth(sheetName, "A", "E", 12)
	xls.SetColWidth(sheetName, "F", "F", 30)
	return xls, nil
}

// CreateStyles generates style ids for all styles used in a report.
func (formatter *ExcelReportFormatter) createStyles(xls *excelize.File) error {

	var err error
	formatter.headlineStyleId, err = xls.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "bottom", Color: "000000", Style: 5},
		},
		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err != nil {
		return err
	}
	formatter.daysBottomStyleId, err = xls.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "bottom", Color: "000000", Style: 5},
		},
	})
	if err != nil {
		return err
	}

	formatter.weekendStyleId, err = xls.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#FEC7CE"}, Pattern: 1},
	})
	if err != nil {
		return err
	}

	formatter.holidayStyleId, err = xls.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#FED7DE"}, Pattern: 1},
	})
	if err != nil {
		return err
	}

	formatter.illnessStyleId, err = xls.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#ff6600"}, Pattern: 1},
	})
	if err != nil {
		return err
	}

	formatter.vacationStyleId, err = xls.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#669900"}, Pattern: 1},
	})
	return err
}

// NewExcelFile creates a new, empty Excel file with one sheet using passed sheet name.
func newExcelFile(sheetName string) *excelize.File {

	xls := excelize.NewFile()
	sheetIdx := 0
	sheets := xls.GetSheetList()
	if len(sheets) > 0 {
		xls.SetSheetName(sheets[0], sheetName)
		sheetIdx = xls.GetSheetIndex(sheetName)
	} else {
		sheetIdx = xls.NewSheet(sheetName)
	}
	xls.SetActiveSheet(sheetIdx)
	return xls
}

// WriteHeader adds headline columns at first row of givem Excel file.
func writeHeader(xls *excelize.File, sheetName string) {
	xls.SetCellValue(sheetName, getCellId("A", 1), "Date")
	xls.SetCellValue(sheetName, getCellId("B", 1), "Start")
	xls.SetCellValue(sheetName, getCellId("C", 1), "End")
	xls.SetCellValue(sheetName, getCellId("D", 1), "WorkingTime")
	xls.SetCellValue(sheetName, getCellId("E", 1), "BreakTime")
	xls.SetCellValue(sheetName, getCellId("F", 1), "Comment")
}

// WriteSummary appends total working time at given row.
func writeSummary(xls *excelize.File, sheetName string, row int, totalWorkingTime time.Duration) {
	xls.SetCellValue(sheetName, getCellId("D", row), formatDuration(totalWorkingTime))
}

// AppendRowForDay will write values for a single day to the Excel file.
// This will apply all required styles or weekends or holidayys as well.
func (formatter *ExcelReportFormatter) appendRowForDay(day Day, xls *excelize.File, sheetName string, row int) {

	date := formatter.atTimezone(day.Date.AsTime())
	xls.SetCellValue(sheetName, getCellId("A", row), date.Format(formatter.dateFormat))
	if len(day.Events) > 0 {
		xls.SetCellValue(sheetName, getCellId("B", row), formatter.formatTime(day.Events[0].Timestamp))
	}
	if len(day.Events) > 1 {
		xls.SetCellValue(sheetName, getCellId("C", row), formatter.formatTime(day.Events[len(day.Events)-1].Timestamp))
	}
	xls.SetCellValue(sheetName, getCellId("D", row), formatDuration(day.WorkingTime))
	xls.SetCellValue(sheetName, getCellId("E", row), formatDuration(day.BreakTime))

	if day.Type == VACATION {
		xls.SetCellStyle(sheetName, getCellId("A", row), getCellId("F", row), formatter.vacationStyleId)
		xls.SetCellValue(sheetName, getCellId("F", row), "Vacation")
	}

	if day.Type == ILLNESS {
		xls.SetCellStyle(sheetName, getCellId("A", row), getCellId("F", row), formatter.illnessStyleId)
		xls.SetCellValue(sheetName, getCellId("F", row), "Illness")
	}

	if isWeekend(day.Date.AsTime()) {
		xls.SetCellStyle(sheetName, getCellId("A", row), getCellId("F", row), formatter.weekendStyleId)
	}
	if holiday, ok := formatter.holidays[day.Date]; ok {
		xls.SetCellStyle(sheetName, getCellId("A", row), getCellId("F", row), formatter.holidayStyleId)
		xls.SetCellValue(sheetName, getCellId("F", row), holiday.Description)
	}
}

// AtTimezone returns passed time in a timezone defined for this formatter.
func (formatter *ExcelReportFormatter) atTimezone(t time.Time) time.Time {
	if formatter.timezone != nil {
		return t.In(formatter.timezone)
	}
	return t
}

// GetCellId helper to generate a cell id by goven column and row index.
func getCellId(column string, row int) string {
	return fmt.Sprintf("%s%d", column, row)
}

// DetermineDateFormat will apply date format if it has been defined in report locale.
// Default format is "2006-01-02".
func (formatter *ExcelReportFormatter) determineDateFormat(report *MonthlyReport) {
	if report.Location.DateFormat != nil {
		formatter.dateFormat = *report.Location.DateFormat
	}
}

// DetermineTimezone will assign a timezone to a formatter it it has been defined in report local.
func (formatter *ExcelReportFormatter) determineTimezone(report *MonthlyReport) {
	if report.Location.Timezone != nil {
		if location, err := time.LoadLocation(*report.Location.Timezone); err == nil {
			formatter.timezone = location
			return
		}
	}
	formatter.timezone = nil
}

// FormatTime applies default time format to given timestamp.
func (formatter *ExcelReportFormatter) formatTime(t time.Time) string {
	return formatter.atTimezone(t).Format(formatter.timeFormat)
}

// FormatDuration returns string representation of given duration in format HH:MM.
func formatDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d:%02d", h, m)
}

// GenerateIndexMap creates a map where date of a day in used as index.
func generateIndexMap(days []Day) map[Date]Day {
	daysMao := make(map[Date]Day)
	for _, day := range days {
		daysMao[day.Date] = day
	}
	return daysMao
}

// EmptyDay returns a working day without time tracking events an no working or break time.
func emptyDay(date Date) Day {
	return Day{
		Date:        date,
		Type:        WORKDAY,
		WorkingTime: time.Duration(0),
		BreakTime:   time.Duration(0),
		Events:      []TimeTrackingRecord{},
	}
}

// AsHolidayMap generates a map with date index for passed lost pf holidays.
func asHolidayMap(holidays []Holiday) map[Date]Holiday {
	holidayMap := make(map[Date]Holiday)
	for _, holiday := range holidays {
		holidayMap[holiday.Date] = holiday
	}
	return holidayMap
}
