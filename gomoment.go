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

// Create Moment Instance
func NewMoment(t time.Time) *Moment {
	return &Moment{t: t}
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
