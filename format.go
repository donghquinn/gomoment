package gomoment

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// formatTokens is the set of supported format tokens. A token is always a run
// of one repeated letter, so "MMM" is rejected rather than read as "MM" + "M".
var formatTokens = map[string]struct{}{
	"YYYY": {}, "YY": {},
	"MM": {}, "M": {},
	"DD": {}, "D": {},
	"HH": {}, "H": {}, "hh": {}, "h": {},
	"mm": {}, "m": {},
	"ss": {}, "s": {}, "SSS": {},
	"A": {}, "a": {},
	"ZZ": {}, "Z": {},
}

// supportedTokens is the token list shown in error messages.
const supportedTokens = "YYYY, YY, MM, M, DD, D, HH, H, hh, h, mm, m, ss, s, SSS, A, a, ZZ, Z"

// Format returns a formatted time string using moment.js-like format tokens.
//
// Supported tokens: YYYY, YY, MM, M, DD, D, HH, H, hh, h, mm, m, ss, s, SSS, A, a, ZZ, Z.
// Text wrapped in square brackets is emitted literally (e.g. "[Today is] YYYY").
// Non-letter characters are copied as-is; unknown ASCII letters cause an error.
func (m *Moment) Format(format string) (string, error) {
	buf := make([]byte, 0, len(format)+16)
	var invalid []string
	invalidEnd := -1

	for i := 0; i < len(format); {
		c := format[i]

		// Escaped literal text: [ ... ]
		if c == '[' {
			if end := strings.IndexByte(format[i+1:], ']'); end >= 0 {
				buf = append(buf, format[i+1:i+1+end]...)
				i += end + 2
				continue
			}
		}

		if !isASCIILetter(c) {
			buf = append(buf, c)
			i++
			continue
		}

		start := i
		for i < len(format) && format[i] == c {
			i++
		}
		token := format[start:i]

		if _, ok := formatTokens[token]; ok {
			buf = appendToken(buf, m.t, token)
			continue
		}

		// Merge adjacent unsupported runs so "XYZ" is reported as one token.
		if n := len(invalid); n > 0 && invalidEnd == start {
			invalid[n-1] += token
		} else {
			invalid = append(invalid, token)
		}
		invalidEnd = i
	}

	if len(invalid) > 0 {
		return "", fmt.Errorf("unsupported format tokens: %v (supported tokens: %s)", invalid, supportedTokens)
	}
	return string(buf), nil
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

// appendToken appends the value of a single format token for t to buf.
func appendToken(buf []byte, t time.Time, token string) []byte {
	switch token {
	case "YYYY":
		return appendInt(buf, t.Year(), 4)
	case "YY":
		return appendInt(buf, t.Year()%100, 2)
	case "MM":
		return appendInt(buf, int(t.Month()), 2)
	case "M":
		return appendInt(buf, int(t.Month()), 1)
	case "DD":
		return appendInt(buf, t.Day(), 2)
	case "D":
		return appendInt(buf, t.Day(), 1)
	case "HH":
		return appendInt(buf, t.Hour(), 2)
	case "H":
		return appendInt(buf, t.Hour(), 1)
	case "hh":
		return appendInt(buf, hour12(t), 2)
	case "h":
		return appendInt(buf, hour12(t), 1)
	case "mm":
		return appendInt(buf, t.Minute(), 2)
	case "m":
		return appendInt(buf, t.Minute(), 1)
	case "ss":
		return appendInt(buf, t.Second(), 2)
	case "s":
		return appendInt(buf, t.Second(), 1)
	case "SSS":
		return appendInt(buf, t.Nanosecond()/int(time.Millisecond), 3)
	case "A":
		if t.Hour() < 12 {
			return append(buf, "AM"...)
		}
		return append(buf, "PM"...)
	case "a":
		if t.Hour() < 12 {
			return append(buf, "am"...)
		}
		return append(buf, "pm"...)
	case "ZZ":
		return appendOffset(buf, t, false)
	case "Z":
		return appendOffset(buf, t, true)
	}
	return buf
}

// hour12 returns the hour in 12-hour clock format (1-12).
func hour12(t time.Time) int {
	h := t.Hour() % 12
	if h == 0 {
		return 12
	}
	return h
}

// appendOffset appends the UTC offset of t as +hhmm, or +hh:mm when colon is true.
func appendOffset(buf []byte, t time.Time, colon bool) []byte {
	_, offset := t.Zone()
	if offset < 0 {
		buf = append(buf, '-')
		offset = -offset
	} else {
		buf = append(buf, '+')
	}
	buf = appendInt(buf, offset/3600, 2)
	if colon {
		buf = append(buf, ':')
	}
	return appendInt(buf, offset%3600/60, 2)
}

// appendInt appends n to buf, left-padded with zeros to at least width digits.
func appendInt(buf []byte, n, width int) []byte {
	if n < 0 {
		buf = append(buf, '-')
		n = -n
	}
	digits := 1
	for v := n; v >= 10; v /= 10 {
		digits++
	}
	for ; digits < width; digits++ {
		buf = append(buf, '0')
	}
	return strconv.AppendInt(buf, int64(n), 10)
}

func isASCIILetter(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}
