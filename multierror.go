package venom

import (
	"reflect"

	"github.com/hashicorp/go-multierror"
)

func allNil(errs []error) bool {
	for _, e := range errs {
		if e != nil {
			return false
		}
		if !reflect.ValueOf(e).IsNil() {
			return false
		}
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
