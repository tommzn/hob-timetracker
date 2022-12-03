package timetracker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type ExcelReportFormatterTestSuite struct {
	suite.Suite
}

func TestExcelReportFormatterTestSuitee(t *testing.T) {
	suite.Run(t, new(ExcelReportFormatterTestSuite))
}

func (suite *ExcelReportFormatterTestSuite) TestGenerateReport() {

	formater := NewExcelReportFormatter()
	report := monthlyReportForTest()

	api, ok := holidayApiForTest()
	suite.True(ok)

	holidays, err := api.GetHolidays(report.Year, report.Month)
	suite.Nil(err)
	formater.WithHolidays(holidays)

	suite.Nil(formater.WriteMonthlyReportToFile(report, "report.xlsx"))
	buf, err := formater.WriteMonthlyReportToBuffer(report)
	suite.Nil(err)
	suite.True(len(buf.Bytes()) > 0)
}

func monthlyReportForTest() *MonthlyReport {
	return &MonthlyReport{
		Year:     2022,
		Month:    1,
		Location: localeForTest(),
		Days: []Day{
			Day{
				Date:        Date{Year: 2022, Month: 1, Day: 1},
				Type:        WORKDAY,
				WorkingTime: 8 * time.Hour,
				BreakTime:   30 * time.Minute,
				Events: []TimeTrackingRecord{
					TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-01-01T08:00:00")},
					TimeTrackingRecord{Type: WORKDAY, Timestamp: asTime("2022-01-01T16:30:00")},
				},
			},
			Day{
				Date:        Date{Year: 2022, Month: 1, Day: 10},
				Type:        ILLNESS,
				WorkingTime: 8 * time.Hour,
				BreakTime:   30 * time.Minute,
				Events:      []TimeTrackingRecord{},
			},
			Day{
				Date:        Date{Year: 2022, Month: 1, Day: 11},
				Type:        VACATION,
				WorkingTime: 8 * time.Hour,
				BreakTime:   30 * time.Minute,
				Events:      []TimeTrackingRecord{},
			},
			Day{
				Date:        Date{Year: 2022, Month: 1, Day: 12},
				Type:        VACATION,
				WorkingTime: 8 * time.Hour,
				BreakTime:   30 * time.Minute,
				Events:      []TimeTrackingRecord{},
			},
		},
		TotalWorkingTime: 8 * time.Hour,
	}
}
