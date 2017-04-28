//
// Here we implement extra flags to extend pflag/viper binding.
//
// https://godoc.org/github.com/spf13/pflag#Value
//
// type Value interface {
//     String() string
//     Set(string) error
//     Type() string
// }
//

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
	val, err := time.Parse(DefaultTimeFormat, s)
	*v = timeValue(val)
	return err
}

func (v *timeValue) Type() string {
	return "time"
}

func (v *timeValue) String() string {
	return time.Time(*v).Format(DefaultTimeFormat)
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
	return "duration"
}

func (v *durationValue) String() string {
	return time.Duration(*v).String()
}
