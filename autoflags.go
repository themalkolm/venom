package venom

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
)

func parseTag(tag string) (string, string, string) {
	parts := strings.SplitN(tag, ",", 3)

	// flag: bar, b, Some barness -> flag: bar,b,Some barness
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}

	switch len(parts) {
	case 1:
		// flag: b
		if len(parts[0]) == 1 {
			return "", parts[0], ""
		}
		// flag: bar
		return parts[0], "", ""
	case 2:
		// flag: b,Some barness
		if len(parts[0]) == 1 {
			return "", parts[0], parts[1]
		}
		// flag: bar,b
		if len(parts[1]) == 1 {
			return parts[0], parts[1], ""
		}
		// flag: bar,Some barness
		return parts[0], "", parts[1]
	case 3:
		// flag: bar,b,Some barness
		return parts[0], parts[1], parts[2]
	default:
		return "", "", ""
	}
}

func DefineFlags(config interface{}) *pflag.FlagSet {
	flags, err := NewFlags(config)
	if err != nil {
		panic(err)
	}
	return flags
}

func NewFlags(config interface{}) (*pflag.FlagSet, error) {
	var flags pflag.FlagSet

	//
	// Remove one level of indirection.
	//
	v := reflect.ValueOf(config)
	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}

	//
	// Make sure we end up with a struct.
	//
	if v.Kind() != reflect.Struct {
		return nil, errors.New("Struct or pointer to struct expected")
	}

	//
	// For every struct field create a flag.
	//
	for i := 0; i < v.Type().NumField(); i++ {
		field := v.Type().Field(i)

		tag := field.Tag.Get("pflag")
		if tag == "" {
			continue
		}

		name, zero, usage := parseTag(tag)

		val := v.Field(i)
		typ := val.Type()
		switch typ.Kind() {
		case reflect.Bool:
			if value, err := strconv.ParseBool(zero); err != nil {
				return nil, err
			} else {
				flags.Bool(name, value, usage)
			}
		case reflect.Int:
			if value, err := strconv.ParseInt(zero, 10, 64); err != nil {
				return nil, err
			} else {
				flags.Int(name, int(value), usage)
			}
		case reflect.Int8:
			if value, err := strconv.ParseInt(zero, 10, 8); err != nil {
				return nil, err
			} else {
				flags.Int8(name, int8(value), usage)
			}
		case reflect.Int16:
			if value, err := strconv.ParseInt(zero, 10, 16); err != nil {
				return nil, err
			} else {
				flags.Int32(name, int32(value), usage) // Not a typo, pflags doesn't have Int16
			}
		case reflect.Int32:
			if value, err := strconv.ParseInt(zero, 10, 32); err != nil {
				return nil, err
			} else {
				flags.Int32(name, int32(value), usage)
			}
		case reflect.Int64:
			if value, err := strconv.ParseInt(zero, 10, 64); err != nil {
				return nil, err
			} else {
				flags.Int64(name, int64(value), usage)
			}
		case reflect.Uint:
			if value, err := strconv.ParseUint(zero, 10, 64); err != nil {
				return nil, err
			} else {
				flags.Uint(name, uint(value), usage)
			}
		case reflect.Uint8:
			if value, err := strconv.ParseUint(zero, 10, 8); err != nil {
				return nil, err
			} else {
				flags.Uint8(name, uint8(value), usage)
			}
		case reflect.Uint16:
			if value, err := strconv.ParseUint(zero, 10, 16); err != nil {
				return nil, err
			} else {
				flags.Uint16(name, uint16(value), usage)
			}
		case reflect.Uint32:
			if value, err := strconv.ParseUint(zero, 10, 32); err != nil {
				return nil, err
			} else {
				flags.Uint32(name, uint32(value), usage)
			}
		case reflect.Uint64:
			if value, err := strconv.ParseUint(zero, 10, 64); err != nil {
				return nil, err
			} else {
				flags.Uint64(name, uint64(value), usage)
			}
		case reflect.Float32:
			if value, err := strconv.ParseFloat(zero, 32); err != nil {
				return nil, err
			} else {
				flags.Float32(name, float32(value), usage)
			}
		case reflect.Float64:
			if value, err := strconv.ParseFloat(zero, 64); err != nil {
				return nil, err
			} else {
				flags.Float64(name, float64(value), usage)
			}
		case reflect.String:
			flags.String(name, zero, usage)
		default:
			return nil, fmt.Errorf("Unsupported type for field with flag tag %q: %s", name, typ)
		}
	}

	return &flags, nil
}
