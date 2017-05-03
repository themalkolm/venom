package venom

import (
	"reflect"

	"github.com/hashicorp/go-multierror"
)

func isNil(e error) bool {
	return e == nil || reflect.ValueOf(e).IsNil()
}

func allNil(errs []error) bool {
	for _, e := range errs {
		if isNil(e) {
			continue
		}

		return false
	}
	return true
}

//
// Same as multierror.Append but takes extra care to not create
// non-nil multierror.Error object with no errors.
//
func AppendErr(err error, errs ...error) error {
	ret := multierror.Append(err, errs...)
	if len(ret.Errors) == 0 || allNil(ret.Errors) {
		return nil
	}
	return ret.ErrorOrNil()
}
