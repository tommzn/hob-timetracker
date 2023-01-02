package timetracker

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

// Format returns a date as string in passed layout.
func (d Date) Format(layout string) string {
	formattedDate := []rune(layout)
	formattedDate = replaceInfFormat(layout, "2006", strconv.Itoa(d.Year), formattedDate)
	formattedDate = replaceInfFormat(layout, "01", fmt.Sprintf("%02d", d.Month), formattedDate)
	formattedDate = replaceInfFormat(layout, "02", fmt.Sprintf("%02d", d.Day), formattedDate)
	return string(formattedDate)
}

// IsEqual return true if given date matched with current.
func (d Date) IsEqual(t time.Time) bool {
	return d.Year == t.Year() && d.Month == int(t.Month()) && d.Day == t.Day()
}

// AsTime returns a time object, at 00:00:00 UtC, for current date.
func (d Date) AsTime() time.Time {
	return time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, time.UTC)
}

// Before returns true if passed date is before current date.
func (d Date) Before(d2 Date) bool {
	return d.Year*10000+d.Month*100+d.Day < d2.Year*10000+d2.Month*100+d2.Day
}

// String returns current date formatted with YYYY-MM-DD/2006-01-02.
func (d Date) String() string {
	return d.Format("2006-01-02")
}

// ReplaceInfFormat is a formatting helper to replace current date values in given format.
// Uses format constants from Golang's time package.
func replaceInfFormat(layout, formatKey, newValue string, target []rune) []rune {
	formatted := []rune(target)
	if len(formatKey) == len(newValue) {
		if idx := strings.Index(layout, formatKey); idx >= 0 {
			for _, c := range []rune(newValue) {
				formatted[idx] = c
				idx++
			}
		}
	}
	return formatted
}

// AsDate is a helper to convert given time to a single date. UTC timestamp of given time is used.
func asDate(t time.Time) Date {
	year, month, day := t.UTC().Date()
	return Date{Year: year, Month: int(month), Day: day}
}

// NextDay increases given date by one day.
func nextDay(day time.Time) time.Time {
	return day.AddDate(0, 0, 1)
}

// IsDayBeforeOrEqual returns with true if first date is less or equal than second date.
func isDayBeforeOrEqual(t1 time.Time, t2 time.Time) bool {
	m1 := time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.UTC)
	m2 := time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, time.UTC)
	return m1.Before(m2) || m1.Equal(m2)
}

// IsWeekend returns true if passed weekday is Saturday or Sunday.
func isWeekend(day time.Time) bool {
	return slices.Contains([]time.Weekday{time.Saturday, time.Sunday}, day.Weekday())
}

// IsInRange returns true if passed timestamp is greater or quals than start
// and less or equal as end.
func isInRange(start, end, timestamp time.Time) bool {
	return (start.UTC().Before(timestamp.UTC()) || start.UTC().Equal(timestamp.UTC())) &&
		(end.UTC().After(timestamp.UTC()) || end.UTC().Equal(timestamp.UTC()))
}
