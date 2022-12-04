package timetracker

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type ReportCalulatorTestSuite struct {
	suite.Suite
}

func TestReportCalulatorTestSuite(t *testing.T) {
	suite.Run(t, new(ReportCalulatorTestSuite))
}

func (suite *ReportCalulatorTestSuite) TestCalcEightHourWorkingDay() {

	records := []TimeTrackingRecord{
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-02-01T08:00:00")},
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-02-01T16:00:00")},
	}
	calculator := NewReportCalulator(records, localeForTest())
	report, err := calculator.MonthlyReport(2022, 2, WORKDAY)
	suite.Nil(err)
	suite.assertWorkingTime(report, 7*time.Hour+30*time.Minute)
}

func (suite *ReportCalulatorTestSuite) TestCalcTenHourWorkingDay() {

	records := []TimeTrackingRecord{
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-02-01T08:00:00")},
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-02-01T18:00:00")},
	}
	calculator := NewReportCalulator(records, localeForTest())
	report, err := calculator.MonthlyReport(2022, 2, WORKDAY)
	suite.Nil(err)
	suite.assertWorkingTime(report, 9*time.Hour+15*time.Minute)
}

func (suite *ReportCalulatorTestSuite) TestCalcShortWorkingDay() {

	records := []TimeTrackingRecord{
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-02-01T08:00:00")},
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-02-01T12:00:00")},
	}
	calculator := NewReportCalulator(records, localeForTest())
	report, err := calculator.MonthlyReport(2022, 2, WORKDAY)
	suite.Nil(err)
	suite.assertWorkingTime(report, 4*time.Hour)
}

func (suite *ReportCalulatorTestSuite) TestCalcWorkingDayWithBreaks() {

	records := []TimeTrackingRecord{
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-02-01T08:00:00")},
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-02-01T10:00:00")},
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-02-01T11:00:00")},
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-02-01T13:00:00")},
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-02-01T14:00:00")},
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-02-01T16:00:00")},
	}
	calculator := NewReportCalulator(records, localeForTest())
	report, err := calculator.MonthlyReport(2022, 2, WORKDAY)
	suite.Nil(err)
	suite.assertWorkingTime(report, 6*time.Hour)
}

func (suite *ReportCalulatorTestSuite) TestCalcWorkingDayWithLongBreak() {

	records := []TimeTrackingRecord{
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-02-01T08:00:00")},
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-02-01T12:00:00")},
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-02-01T13:30:00")},
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-02-01T18:30:00")},
	}
	calculator := NewReportCalulator(records, localeForTest())
	report, err := calculator.MonthlyReport(2022, 2, WORKDAY)
	suite.Nil(err)
	suite.assertWorkingTime(report, 9*time.Hour)
}

func (suite *ReportCalulatorTestSuite) TestIgnoreOtherMonthy() {

	records := []TimeTrackingRecord{
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-01-30T08:00:00")},
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-01-30T12:00:00")},
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-02-01T08:15:00")},
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-02-01T18:00:00")},
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-03-01T08:00:00")},
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-03-01T18:30:00")},
	}
	calculator := NewReportCalulator(records, localeForTest())
	report, err := calculator.MonthlyReport(2022, 2, WORKDAY)
	suite.Nil(err)
	suite.assertWorkingTime(report, 9*time.Hour)
}

func (suite *ReportCalulatorTestSuite) TestFillVacationAndIllness() {

	records := []TimeTrackingRecord{
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-02-07T98:00:00")},
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-02-07T17:00:00")},
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-02-08T98:00:00")},
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-02-08T17:00:00")},
		TimeTrackingRecord{Type: ILLNESS, Timestamp: asTime("2022-02-10T08:00:00")},
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-02-14T98:00:00")},
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-02-14T17:00:00")},
	}
	calculator := NewReportCalulator(records, localeForTest())
	report, err := calculator.MonthlyReport(2022, 2, VACATION)
	suite.Nil(err)
	suite.NotNil(report)
	suite.Len(report.Days, 14)

	for _, day := range report.Days {
		suite.Equal(2022, day.Date.Year)
		suite.Equal(2, day.Date.Month)
		switch true {
		case day.Day >= 1 && day.Day <= 6:
			suite.Equal(VACATION, day.Type)
		case day.Day >= 7 && day.Day <= 8:
			suite.Equal(WORKDAY, day.Type)
		case day.Day >= 9 && day.Day <= 13:
			suite.Equal(ILLNESS, day.Type)
		case day.Day >= 14:
			suite.Equal(WORKDAY, day.Type)
		}
	}
}

func (suite *ReportCalulatorTestSuite) TestCalcEndOfWorkingDay() {

	records1 := []TimeTrackingRecord{
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-02-01T08:00:00")},
	}
	calculator1 := NewReportCalulator(records1, localeForTest())
	report1, err1 := calculator1.MonthlyReport(2022, 2, WORKDAY)
	suite.Nil(err1)
	suite.assertWorkingTime(report1, 8*time.Hour)

	records2 := []TimeTrackingRecord{
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-02-01T08:00:00")},
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-02-01T18:00:00")},
		TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-02-01T18:20:00")},
	}
	calculator2 := NewReportCalulator(records2, localeForTest())
	report2, err2 := calculator2.MonthlyReport(2022, 2, WORKDAY)
	suite.Nil(err2)
	suite.assertWorkingTime(report2, 9*time.Hour+36*time.Minute)
}

func (suite *ReportCalulatorTestSuite) TestDetermineTypeOfDay() {

	emtpyListOfRecords := []TimeTrackingRecord{}
	calculator := NewReportCalulator(emtpyListOfRecords, localeForTest())

	suite.Equal(WORKDAY, calculator.determineTypeOf(emtpyListOfRecords))
}

func (suite *ReportCalulatorTestSuite) TestCalulateWorkTime() {

	emtpyListOfRecords := []TimeTrackingRecord{}
	calculator := NewReportCalulator(emtpyListOfRecords, localeForTest())

	day := Day{
		Date:        asDate(time.Now()),
		Type:        WORKDAY,
		WorkingTime: time.Duration(1 * time.Hour),
		BreakTime:   time.Duration(1 * time.Hour),
		Events:      []TimeTrackingRecord{},
	}
	calculator.calculateWorkTimeForDay(&day)
	suite.Equal(time.Duration(0), day.WorkingTime)
	suite.Equal(time.Duration(0), day.BreakTime)
}

func (suite *ReportCalulatorTestSuite) assertWorkingTime(report *MonthlyReport, expexctedWorkingTime time.Duration) {
	suite.NotNil(report)
	suite.Len(report.Days, 1)
	suite.Equal(expexctedWorkingTime, report.Days[0].WorkingTime)
	suite.Equal(WORKDAY, report.Days[0].Type)
}
