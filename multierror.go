package venom

import (
	"reflect"

	"github.com/hashicorp/go-multierror"
)

func isNil(e error) bool {
	if e == nil {
		return true
	}

	// structs can't be nil
	if reflect.TypeOf(e).Kind() == reflect.Struct {
		return false
	}

	return reflect.ValueOf(e).IsNil()
}

func anyNonNil(errs []error) bool {
	for _, e := range errs {
		if !isNil(e) {
			return true
		}
	}
	return false
}

//
// Same as multierror.Append but takes extra care to not create non-nil multierror.Error object with no errors.n
//
// multierror.Append(nil, nil) -> &Error{}
// AppendErr(nil, nil) -> nil
//
func AppendErr(err error, errs ...error) error {
	ret := multierror.Append(err, errs...)
	if len(ret.Errors) == 0 || !anyNonNil(ret.Errors) {
		return nil
	}
	return ret.ErrorOrNil()
}
