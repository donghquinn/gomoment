package gomoment

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Moment wraps time.Time to provide moment.js-like functionality
type Moment struct {
	t time.Time
}

// Package-level variables for performance optimization
var (
	// replacements maps moment.js tokens to Go time format tokens
	replacements = map[string]string{
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

	// sortedTokens holds tokens sorted by length (longest first) to prevent replacement conflicts
	sortedTokens []string

	// supportedTokens holds all supported format tokens for validation
	supportedTokens []string

	// alphaRegex is pre-compiled regex for finding alphabetic tokens
	alphaRegex = regexp.MustCompile("[a-zA-Z]+")

	// dateFormats lists supported input date formats for parsing
	dateFormats = []string{
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
)

// init initializes package-level variables for optimal performance
func init() {
	// Create sorted tokens slice
	sortedTokens = make([]string, 0, len(replacements))
	supportedTokens = make([]string, 0, len(replacements))

	for token := range replacements {
		sortedTokens = append(sortedTokens, token)
		supportedTokens = append(supportedTokens, token)
	}

	// Sort by length (longest first) to prevent replacement conflicts
	sort.Slice(sortedTokens, func(i, j int) bool {
		return len(sortedTokens[i]) > len(sortedTokens[j])
	})

	// Sort supported tokens for consistent error messages
	sort.Strings(supportedTokens)
}

// Format returns a formatted time string using moment.js-like format tokens.
// Supported tokens: YYYY, YY, MM, M, DD, D, HH, H, hh, h, mm, m, ss, s, SSS, A, a, ZZ, Z
func (m *Moment) Format(format string) (string, error) {
	// Validate input format tokens
	tmpFormat := format
	for _, token := range supportedTokens {
		tmpFormat = strings.ReplaceAll(tmpFormat, token, "")
	}

	// Find any remaining alphabetic tokens (invalid)
	invalidTokens := alphaRegex.FindAllString(tmpFormat, -1)
	if len(invalidTokens) > 0 {
		return "", fmt.Errorf("unsupported format tokens: %v (supported tokens: %v)", invalidTokens, supportedTokens)
	}

	// Replace moment.js tokens with Go time format tokens
	result := format
	for _, token := range sortedTokens {
		result = strings.ReplaceAll(result, token, replacements[token])
	}

	return m.t.Format(result), nil
}

// Time returns the underlying time.Time value
func (m *Moment) Time() time.Time {
	return m.t
}

// NewMoment creates a new Moment instance.
// With no arguments, returns current time.
// With string argument, parses the date string.
// With time.Time argument, wraps the time value.
func NewMoment(args ...interface{}) (*Moment, error) {
	if len(args) == 0 {
		// Return current time without any argument
		return &Moment{t: time.Now()}, nil
	}

	switch v := args[0].(type) {
	case string:
		// Parse string
		return Parse(v)
	case time.Time:
		// Wrap time.Time
		return &Moment{t: v}, nil
	default:
		return nil, fmt.Errorf("unsupported argument type: %T (supported: string, time.Time, or no arguments)", args[0])
	}
}

// Parse creates a Moment from a date string.
// Supports various common date formats including ISO 8601.
func Parse(dateStr string) (*Moment, error) {
	// Try parsing with UTC first
	for _, layout := range dateFormats {
		t, err := time.Parse(layout, dateStr)
		if err == nil {
			return &Moment{t: t}, nil
		}
	}

	// Retry with local timezone
	for _, layout := range dateFormats {
		t, err := time.ParseInLocation(layout, dateStr, time.Local)
		if err == nil {
			return &Moment{t: t}, nil
		}
	}

	return nil, fmt.Errorf(`unable to parse date string "%s" (supported formats: YYYY-MM-DD, YYYY/MM/DD, RFC3339, etc.)`, dateStr)
}

// Now creates a new Moment instance with the current time
func Now() *Moment {
	return &Moment{t: time.Now()}
}

// Must formats the time using the given format string and panics on error.
// Use this when you are certain the format string is valid.
func (m *Moment) Must(format string) string {
	result, err := m.Format(format)
	if err != nil {
		panic(err)
	}
	return result
}
