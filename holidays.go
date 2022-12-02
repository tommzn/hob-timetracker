package timetracker

import (
	"strings"

	"github.com/calendarific/go-calendarific"
)

// newCalendarApi returns a new api to get holidays.
func newCalendarApi(apiKey string, location Locale) *CalendarApi {
	return &CalendarApi{apiKey: apiKey, location: location}
}

// CalendarApi is used to retrieve holidays from an online service.
type CalendarApi struct {
	apiKey   string
	location Locale
}

// GetHolidays try to fetch holidays for given month.
func (api *CalendarApi) GetHolidays(year, month int) (map[Date]string, error) {

	listOfHolidays := make(map[Date]string)
	calParams := calendarific.CalParameters{
		ApiKey:  api.apiKey,
		Country: api.location.Country,
		Year:    int32(year),
		Month:   int32(month),
	}

	res, err := calParams.CalData()
	if err != nil {
		return listOfHolidays, err
	}
	for _, holiday := range res.Response.Holidays {
		date := Date{holiday.Date.Datetime.Year, holiday.Date.Datetime.Month, holiday.Date.Datetime.Day}
		if isHoliday(holiday.Type) {
			listOfHolidays[date] = holiday.Name
		}
	}
	return listOfHolidays, nil
}

// IsHoliday determines if given holiday types contains works "holiday" or "national".
func isHoliday(types []string) bool {
	for _, dayType := range types {
		if strings.Contains(strings.ToLower(dayType), "holiday") ||
			strings.Contains(strings.ToLower(dayType), "national") {
			return true
		}
	}
	return false
}
