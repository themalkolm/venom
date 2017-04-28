package venom

import (
	"time"
)

var (
	DefaultTimeFormat = time.RFC3339
)

//
// time.Time
//
type timeValue time.Time

func newTimeValue(val time.Time, p *time.Time) *timeValue {
	*p = val
	return (*timeValue)(p)
}

func (v *timeValue) Set(s string) error {
	val, err := time.Parse(s, DefaultTimeFormat)
	*v = timeValue(val)
	return err
}

func (v *timeValue) Type() string {
	return "time.Time"
}

func (v *timeValue) String() string {
	return time.Time(*v).String()
}

//
// time.Duration
//
type durationValue time.Duration

func newDurationValue(val time.Duration, p *time.Duration) *durationValue {
	*p = val
	return (*durationValue)(p)
}

func (v *durationValue) Set(s string) error {
	val, err := time.ParseDuration(s)
	*v = durationValue(val)
	return err
}

func (v *durationValue) Type() string {
	return "time.Duration"
}

func (v *durationValue) String() string {
	return time.Duration(*v).String()
}
