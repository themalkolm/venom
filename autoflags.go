package venom

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/pflag"
)

const (
	SquashFlagsTag = "++"
)

//
// Structures implementing this interface won't be introspected and this function will be called
// instead.
//
type HasFlags interface {
	Flags() *pflag.FlagSet
}

// Parse name for mapstructure tags i.e. fetch banana from:
//
// type Foo struct {
//     foo int `mapstructure:"banana"`
// }
//
func parseMapstructureTag(tag string) (string, bool) {
	parts := strings.SplitN(tag, ",", 2)
	if len(parts) == 0 {
		return "", false
	}
	return parts[0], true
}

type flagInfo struct {
	name      string
	shorthand string
	usage     string
}

func parseTag(tag string) flagInfo {
	parts := strings.SplitN(tag, ",", 3)

	// flag: bar, b, Some barness -> flag: bar,b,Some barness
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}

	var f flagInfo
	switch len(parts) {
	case 1:
		// flag: b
		if len(parts[0]) == 1 {
			f.name = ""
			f.shorthand = parts[0]
			f.usage = ""
			return f
		}
		// flag: bar
		f.name = parts[0]
		f.shorthand = ""
		f.usage = ""
		return f
	case 2:
		// flag: b,Some barness
		if len(parts[0]) == 1 {
			f.name = ""
			f.shorthand = parts[0]
			f.usage = parts[1]
			return f
		}
		// flag: bar,b
		if len(parts[1]) == 1 {
			f.name = parts[0]
			f.shorthand = parts[1]
			f.usage = ""
			return f
		}
		// flag: bar,Some barness
		f.name = parts[0]
		f.shorthand = ""
		f.usage = parts[1]
		return f
	case 3:
		// flag: bar,b,Some barness
		f.name = parts[0]
		f.shorthand = parts[1]
		f.usage = parts[2]
		return f
	default:
		return f
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

func (a flagsFactory) lookupTag(tag reflect.StructTag) (string, bool) {
	for _, name := range a.tags {
		v, ok := tag.Lookup(name)
		if ok {
			return v, true
		}
	}
	return "", false
}

func (a flagsFactory) createFlags(defaults interface{}) (*pflag.FlagSet, error) {
	var flags pflag.FlagSet

	//
	// Remove one level of indirection.
	//
	v := reflect.ValueOf(defaults)
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
	// For every tagged struct field create a flag.
	//
	for i := 0; i < v.Type().NumField(); i++ {
		structField := v.Type().Field(i)
		fieldType := structField.Type
		fieldValue := v.Field(i)

		tag, ok := a.lookupTag(structField.Tag)
		if !ok {
			continue
		}

		//
		// This means we want to squash all flags from the inner structure so they appear as is they are defined
		// in the outer structure.
		//
		if tag == SquashFlagsTag {
			if fieldType.Kind() != reflect.Struct {
				return nil, fmt.Errorf(`flag:"%s" is supported only for inner structs but is set on: %s`, SquashFlagsTag, fieldType)
			}

			// Check if the struct implements HasFlags right away
			if hasFlags, ok := fieldValue.Interface().(HasFlags); ok {
				innerFlags := hasFlags.Flags()
				flags.AddFlagSet(innerFlags)
				continue
			}

			// Check if struct-ptr implements HasFlags
			if fieldValue.CanAddr() {
				fieldValueAddr := fieldValue.Addr()

				if hasFlags, ok := fieldValueAddr.Interface().(HasFlags); ok {
					innerFlags := hasFlags.Flags()
					flags.AddFlagSet(innerFlags)
					continue
				}
			}

			// Check if inner struct implements HasFlags.
			//
			// I can't manage to get a pointer to inner struct here, it is not addressable and etc. Just as a workaround
			// we make a temporary copy and get a pointer to it instead. Suboptimal but meh, config struct are supposed
			// to be cheap to copy. Note that fieldValueCopy is a pointer.
			//
			fieldValueCopy := reflect.New(fieldType)
			fieldValueCopy.Elem().Set(fieldValue)

			if hasFlags, ok := fieldValueCopy.Interface().(HasFlags); ok {
				innerFlags := hasFlags.Flags()
				flags.AddFlagSet(innerFlags)
				continue
			}

			// No overrides are provided, continue with recursive introspection
			innerFlags, err := a.createFlags(fieldValue.Interface())
			if err != nil {
				return nil, err
			}
			flags.AddFlagSet(innerFlags)
			continue
		}

		//
		// In case we have mapstructure defined it must have exactly the same name as flag has.
		//
		mapTag, ok := structField.Tag.Lookup("mapstructure")
		if ok {
			mapName, ok := parseMapstructureTag(mapTag)
			if ok && !(tag == mapName || strings.HasPrefix(tag, mapName+",")) {
				return nil, fmt.Errorf(`Both "mapstructure" and "flag" tags must have equal names but are different for field: %s`, structField.Name)
				continue
			}
		}

		fi := parseTag(tag)
		err := addFlagForTag(&flags, fi, fieldValue, fieldType)
		if err != nil {
			return nil, err
		}
	}

	return &flags, nil
}

func addFlagForTag(flags *pflag.FlagSet, fi flagInfo, fieldValue reflect.Value, fieldType reflect.Type) error {
	name := fi.name
	shorthand := fi.shorthand
	usage := fi.usage

	switch fieldType.Kind() {
	case reflect.Bool:
		value := bool(fieldValue.Bool())
		flags.BoolP(name, shorthand, value, usage)
	case reflect.Int:
		value := int(fieldValue.Int())
		flags.IntP(name, shorthand, value, usage)
	case reflect.Int8:
		value := int8(fieldValue.Int())
		flags.Int8P(name, shorthand, value, usage)
	case reflect.Int16:
		value := int32(fieldValue.Int())
		flags.Int32P(name, shorthand, value, usage) // Not a typo, pflags doesn't have Int16
	case reflect.Int32:
		value := int32(fieldValue.Int())
		flags.Int32P(name, shorthand, value, usage)
	case reflect.Int64:
		value := int64(fieldValue.Int())
		flags.Int64P(name, shorthand, value, usage)
	case reflect.Uint:
		value := uint(fieldValue.Uint())
		flags.UintP(name, shorthand, value, usage)
	case reflect.Uint8:
		value := uint8(fieldValue.Uint())
		flags.Uint8P(name, shorthand, value, usage)
	case reflect.Uint16:
		value := uint16(fieldValue.Uint())
		flags.Uint16P(name, shorthand, value, usage)
	case reflect.Uint32:
		value := uint32(fieldValue.Uint())
		flags.Uint32P(name, shorthand, value, usage)
	case reflect.Uint64:
		value := uint64(fieldValue.Uint())
		flags.Uint64P(name, shorthand, value, usage)
	case reflect.Float32:
		value := float32(fieldValue.Float())
		flags.Float32P(name, shorthand, value, usage)
	case reflect.Float64:
		value := float64(fieldValue.Float())
		flags.Float64P(name, shorthand, value, usage)
	case reflect.String:
		value := string(fieldValue.String())
		flags.StringP(name, shorthand, value, usage)
	default:
		return fmt.Errorf("Unsupported type for field with flag tag %q: %s", name, fieldType)
	}
	return nil
}
