package gomoment

import (
	"fmt"
	"strings"
	"time"
)

// timezoneAliases maps common timezone abbreviations to IANA location names.
var timezoneAliases = map[string]string{
	"UTC": "UTC",
	"GMT": "GMT",
	"KST": "Asia/Seoul",
	"JST": "Asia/Tokyo",
	"EST": "America/New_York",
	"PST": "America/Los_Angeles",
	"CST": "America/Chicago",
	"MST": "America/Denver",
	"CET": "Europe/Paris",
	"IST": "Asia/Kolkata",
}

// TZ converts the moment to the specified timezone.
// Accepts abbreviations like "UTC", "KST", "EST" (case-insensitive)
// or IANA names like "America/New_York".
func (m *Moment) TZ(timezone string) (*Moment, error) {
	name := timezone
	if fullName, ok := timezoneAliases[strings.ToUpper(timezone)]; ok {
		name = fullName
	}

	loc, err := time.LoadLocation(name)
	if err != nil {
		return nil, fmt.Errorf("invalid timezone: %s", timezone)
	}

	return &Moment{t: m.t.In(loc)}, nil
}

// UTC converts the moment to UTC timezone.
func (m *Moment) UTC() *Moment {
	return &Moment{t: m.t.UTC()}
}

// Local converts the moment to the local timezone.
func (m *Moment) Local() *Moment {
	return &Moment{t: m.t.Local()}
}

// Zone returns the timezone name and offset in seconds east of UTC.
func (m *Moment) Zone() (string, int) {
	return m.t.Zone()
}

// Offset returns the timezone offset in seconds east of UTC.
func (m *Moment) Offset() int {
	_, offset := m.t.Zone()
	return offset
}
