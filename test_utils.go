package timetracker

import (
	"time"

	log "github.com/tommzn/go-log"
)

// loggerForTest creates a new stdout logger for testing.
func loggerForTest() log.Logger {
	return log.NewLogger(log.Debug, nil, nil)
}

// AsStringPointer returns a pointer for given string value.
func asStringPointer(s string) *string {
	return &s
}

// Coverts passed date/time to time. Expects format: 2006-01-02T15:04:05"
func asTime(value string) time.Time {
	t, _ := time.Parse("2006-01-02T15:04:05", value)
	return t
}

// DeviceIdForTest generates an id for testing.
func deviceIdForTest() string {
	return "0cd21f4b-dcc6-4c86-a319-d410b67b6ee0"
}

// LocaleForTest returns a localization for testing.
func localeForTest() Locale {
	return Locale{
		Country:    "de",
		Timezone:   asStringPointer("Europe/Berlin"),
		DateFormat: asStringPointer("02.01.2006"),
		Breaks: map[time.Duration]time.Duration{
			6 * time.Hour: 30 * time.Minute,
			9 * time.Hour: 15 * time.Minute,
		},
		DefaultWorkTime: 8*time.Hour + 30*time.Minute,
	}
}
