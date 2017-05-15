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
	"fmt"
	"strings"
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

//
// map[string]string
//
type mapStringStringValue map[string]string

func newMapStringStringValue(val map[string]string, p map[string]string) mapStringStringValue {
	for k := range val {
		p[k] = val[k]
	}
	return mapStringStringValue(p)
}

func (v mapStringStringValue) Set(s string) error {
	val, err := parseMapStringString(s, ",", "=")
	if err != nil {
		return err
	}

	// clear
	for k := range v {
		delete(v, k)
	}
	// set
	for k := range val {
		v[k] = val[k]
	}
	return nil
}

func (v mapStringStringValue) Type() string {
	return "map[string]string"
}

func (v mapStringStringValue) String() string {
	parts := make([]string, 0, len(v))
	for k := range v {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v[k]))
	}
	return strings.Join(parts, ",")
}
