package gomoment

import (
	"fmt"
	"time"
)

// dateFormats lists supported input date formats for parsing, tried in order.
var dateFormats = []string{
	"2006-01-02",
	"2006/01/02",
	"2006-01-02 15:04:05",
	"2006/01/02 15:04:05",
	"2006-01-02T15:04:05",
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05-07:00",
	"01/02/2006",
	"01-02-2006",
	"01/02/2006 15:04:05",
	"01-02-2006 15:04:05",
	"15:04:05",
}

// Parse creates a Moment from a date string.
// Supports various common date formats including ISO 8601.
// Strings without an explicit offset are interpreted as UTC.
func Parse(dateStr string) (*Moment, error) {
	for _, layout := range dateFormats {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return &Moment{t: t}, nil
		}
	}

	return nil, fmt.Errorf(`unable to parse date string "%s" (supported formats: YYYY-MM-DD, YYYY/MM/DD, RFC3339, etc.)`, dateStr)
}
