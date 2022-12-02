package timetracker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type DateTestSuite struct {
	suite.Suite
}

func TestDateTestSuite(t *testing.T) {
	suite.Run(t, new(DateTestSuite))
}

func (suite *DateTestSuite) TestConversion() {

	t1 := time.Now().UTC()
	d1 := asDate(t1)
	suite.Equal(t1.Year(), d1.Year)
	suite.Equal(int(t1.Month()), d1.Month)
	suite.Equal(t1.Day(), d1.Day)

	layout := "20060102"
	suite.Equal(t1.Format(layout), d1.Format(layout))
}

func (suite *DateTestSuite) TestIsWeekend() {
	suite.True(isWeekend(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)))
	suite.True(isWeekend(time.Date(2022, time.January, 2, 12, 0, 0, 0, time.UTC)))
}
