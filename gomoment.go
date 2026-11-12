// Package gomoment provides moment.js-like date parsing and formatting for Go.
package gomoment

import (
	"fmt"
	"time"
)

// Moment wraps time.Time to provide moment.js-like functionality.
// A Moment is immutable: every conversion returns a new Moment.
type Moment struct {
	t time.Time
}

// NewMoment creates a new Moment instance.
// With no arguments, returns current time.
// With string argument, parses the date string.
// With time.Time argument, wraps the time value.
func NewMoment(args ...interface{}) (*Moment, error) {
	if len(args) == 0 {
		return Now(), nil
	}

	switch v := args[0].(type) {
	case string:
		return Parse(v)
	case time.Time:
		return &Moment{t: v}, nil
	default:
		return nil, fmt.Errorf("unsupported argument type: %T (supported: string, time.Time, or no arguments)", args[0])
	}
}

// Now creates a new Moment instance with the current time.
func Now() *Moment {
	return &Moment{t: time.Now()}
}

// Time returns the underlying time.Time value.
func (m *Moment) Time() time.Time {
	return m.t
}
