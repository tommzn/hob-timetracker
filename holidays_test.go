package timetracker

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type HolidaysTestSuite struct {
	suite.Suite
}

func TestHolidaysTestSuite(t *testing.T) {
	suite.Run(t, new(HolidaysTestSuite))
}

func (suite *HolidaysTestSuite) skipCI() {
	if _, isSet := os.LookupEnv("CI"); isSet {
		suite.T().Skip("Skip test in CI environment.")
	}
}

func (suite *HolidaysTestSuite) TestGetHolidays() {

	suite.skipCI()

	api, ok := holidayApiForTest()
	suite.True(ok)

	holidays, err := api.GetHolidays(2021, 12)
	suite.Nil(err)
	suite.True(len(holidays) > 1)
}

func holidayApiForTest() (*CalendarApi, bool) {
	apiKey, ok := os.LookupEnv("HOLIDAYS_API_KEY")
	return newCalendarApi(apiKey, localeForTest()), ok
}
