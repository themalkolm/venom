package venom

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
)

func stringToTimeDurationHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(time.Duration(5)) {
			return data, nil
		}

		return time.ParseDuration(data.(string))
	}
}

func stringToTimeHookFunc(layout string) mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		// Convert it by parsing
		return time.Parse(layout, data.(string))
	}
}

func stringToStringSliceHookFunc(sep string) mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || (t != reflect.TypeOf([]string{})) {
			return data, nil
		}

		raw := data.(string)
		if raw == "" {
			return []string{}, nil
		}

		return strings.Split(raw, sep), nil
	}
}

func stringToBoolSliceHookFunc(sep string) mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || (t != reflect.TypeOf([]bool{})) {
			return data, nil
		}

		raw := data.(string)
		raw = strings.TrimPrefix(raw, "[")
		raw = strings.TrimSuffix(raw, "]")

		vals := make([]bool, 0)
		if raw == "" {
			return vals, nil
		}

		for _, s := range strings.Split(raw, sep) {
			v, err := strconv.ParseBool(s)
			if err != nil {
				return nil, err
			}
			vals = append(vals, v)
		}
		return vals, nil
	}
}

func stringToIntSliceHookFunc(sep string) mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || (t != reflect.TypeOf([]int{})) {
			return data, nil
		}

		raw := data.(string)
		raw = strings.TrimPrefix(raw, "[")
		raw = strings.TrimSuffix(raw, "]")

		vals := make([]int, 0)
		if raw == "" {
			return vals, nil
		}

		for _, s := range strings.Split(raw, sep) {
			v, err := strconv.ParseInt(s, 10, 0)
			if err != nil {
				return nil, err
			}
			vals = append(vals, int(v))
		}
		return vals, nil
	}
}

func stringToUintSliceHookFunc(sep string) mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || (t != reflect.TypeOf([]uint{})) {
			return data, nil
		}

		raw := data.(string)
		raw = strings.TrimPrefix(raw, "[")
		raw = strings.TrimSuffix(raw, "]")

		vals := make([]uint, 0)
		if raw == "" {
			return vals, nil
		}

		for _, s := range strings.Split(raw, sep) {
			v, err := strconv.ParseUint(s, 10, 0)
			if err != nil {
				return nil, err
			}
			vals = append(vals, uint(v))
		}
		return vals, nil
	}
}
