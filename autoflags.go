package venom

import (
	"errors"
	"fmt"
	"reflect"
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
	a := flagsFactory{
		tags: []string{"flag", "pflag"},
	}
	return a.createFlags(config)
}

type flagsFactory struct {
	tags []string
}

func (a flagsFactory) lookupTag(field reflect.StructField) (string, bool) {
	for _, name := range a.tags {
		v, ok := field.Tag.Lookup(name)
		if ok {
			return v, true
		}
	}
	return "", false
}

func (a flagsFactory) createFlags(config interface{}) (*pflag.FlagSet, error) {
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
		fieldType := v.Type().Field(i)

		tag, ok := a.lookupTag(fieldType)
		if !ok {
			continue
		}

		name, shorthand, usage := parseTag(tag)

		val := v.Field(i)
		typ := val.Type()
		switch typ.Kind() {
		case reflect.Bool:
			flags.BoolP(name, shorthand, false, usage)
		case reflect.Int:
			flags.IntP(name, shorthand, 0, usage)
		case reflect.Int8:
			flags.Int8P(name, shorthand, 0, usage)
		case reflect.Int16:
			flags.Int32P(name, shorthand, 0, usage) // Not a typo, pflags doesn't have Int16
		case reflect.Int32:
			flags.Int32P(name, shorthand, 0, usage)
		case reflect.Int64:
			flags.Int64P(name, shorthand, 0, usage)
		case reflect.Uint:
			flags.UintP(name, shorthand, 0, usage)
		case reflect.Uint8:
			flags.Uint8P(name, shorthand, 0, usage)
		case reflect.Uint16:
			flags.Uint16P(name, shorthand, 0, usage)
		case reflect.Uint32:
			flags.Uint32P(name, shorthand, 0, usage)
		case reflect.Uint64:
			flags.Uint64P(name, shorthand, 0, usage)
		case reflect.Float32:
			flags.Float32P(name, shorthand, 0, usage)
		case reflect.Float64:
			flags.Float64P(name, shorthand, 0, usage)
		case reflect.String:
			flags.StringP(name, shorthand, "", usage)
		default:
			return nil, fmt.Errorf("Unsupported type for field with flag tag %q: %s", name, typ)
		}
	}

	return &flags, nil
}
