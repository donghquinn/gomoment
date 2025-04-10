package gomoment

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Save Time value
type Moment struct {
	t time.Time
}

// Supporting Time formats
var supportedFormats = []string{
	"YYYY", "YY", "MM", "M", "DD", "D",
	"HH", "H", "hh", "h", "mm", "m",
	"ss", "s", "SSS", "A", "a", "ZZ", "Z",
}

// Date Format List for parsing date string
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

// Make Formatted time string according to the given format
func (m *Moment) Format(format string) (string, error) {
	// Same usage with moment.js
	replacements := map[string]string{
		"YYYY": "2006",
		"YY":   "06",
		"MM":   "01",
		"M":    "1",
		"DD":   "02",
		"D":    "2",
		"HH":   "15",
		"H":    "15",
		"hh":   "03",
		"h":    "3",
		"mm":   "04",
		"m":    "4",
		"ss":   "05",
		"s":    "5",
		"SSS":  "000",
		"A":    "PM",
		"a":    "pm",
		"ZZ":   "-0700",
		"Z":    "-07:00",
	}

	tokens := make([]string, 0, len(replacements))
	for token := range replacements {
		tokens = append(tokens, token)
	}

	sort.Slice(tokens, func(i, j int) bool {
		return len(tokens[i]) > len(tokens[j])
	})

	// Validate input formats
	tmpFormat := format
	for _, token := range supportedFormats {
		tmpFormat = strings.ReplaceAll(tmpFormat, token, "")
	}

	alphaRegex := regexp.MustCompile("[a-zA-Z]+")
	invalidTokens := alphaRegex.FindAllString(tmpFormat, -1)

	if len(invalidTokens) > 0 {
		return "", fmt.Errorf("invalid token: %v", invalidTokens)
	}

	result := format
	for _, token := range tokens {
		result = strings.ReplaceAll(result, token, replacements[token])
	}

	return m.t.Format(result), nil
}

func (m *Moment) Time() time.Time {
	return m.t
}

// Create Moment Instance
func NewMoment(args ...interface{}) (*Moment, error) {
	if len(args) == 0 {
		// return current time without any argument
		return &Moment{t: time.Now()}, nil
	}

	switch v := args[0].(type) {
	case string:
		// Pasre string
		return Parse(v)
	case time.Time:
		// time.Time
		return &Moment{t: v}, nil
	default:
		return nil, fmt.Errorf("invalid argument type: %T", args[0])
	}
}

// Parse String in order to create momemnt
func Parse(dateStr string) (*Moment, error) {
	for _, layout := range dateFormats {
		t, err := time.Parse(layout, dateStr)
		if err == nil {
			return &Moment{t: t}, nil
		}
	}

	// Retry with local time
	for _, layout := range dateFormats {
		t, err := time.ParseInLocation(layout, dateStr, time.Local)
		if err == nil {
			return &Moment{t: t}, nil
		}
	}

	return nil, fmt.Errorf("invalid time format: %s", dateStr)
}

// Create Moment Instance with Current Time
func Now() *Moment {
	return &Moment{t: time.Now()}
}

// Verify errors and throw panic if it exists
// or return formatted formatted strings
func (m *Moment) Must(format string) string {
	result, err := m.Format(format)
	if err != nil {
		panic(err)
	}
	return result
}

func (m *Moment) MustFormat(format string) string {
	result, err := m.Format(format)
	if err != nil {
		panic(err)
	}
	return result
}
