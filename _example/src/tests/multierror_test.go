package tests

import (
	"io"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"

	"github.com/themalkolm/venom"
)

type pair struct {
	err  error
	errs []error
}

var (
	tableNil = map[string]pair{
		"nil-nil":       {err: nil, errs: nil},
		"nil-[]":        {err: nil, errs: []error{}},
		"nil-[nil]":     {err: nil, errs: []error{nil}},
		"nil-[nil,nil]": {err: nil, errs: []error{nil, nil}},
		"{}-nil":        {err: &multierror.Error{}, errs: nil},
		"{}-[]":         {err: &multierror.Error{}, errs: []error{}},
		"{}-[nil]":      {err: &multierror.Error{}, errs: []error{nil}},
		"{}-[nil,nil]":  {err: &multierror.Error{}, errs: []error{nil, nil}},
	}
	tableNonNil = map[string]pair{
		"{EOF}-nil":       {err: &multierror.Error{Errors: []error{io.EOF}}, errs: nil},
		"{EOF}-[]":        {err: &multierror.Error{Errors: []error{io.EOF}}, errs: []error{}},
		"{EOF}-[nil]":     {err: &multierror.Error{Errors: []error{io.EOF}}, errs: []error{nil}},
		"{EOF}-[nil,nil]": {err: &multierror.Error{Errors: []error{io.EOF}}, errs: []error{nil, nil}},
	}
)

func TestAppend_Nil(t *testing.T) {
	for name, p := range tableNil {
		t.Run(name, func(t *testing.T) {
			ret := multierror.Append(p.err, p.errs...)
			assert.NotNil(t, ret)
		})
	}
}

func TestAppend_NonNil(t *testing.T) {
	for name, p := range tableNonNil {
		t.Run(name, func(t *testing.T) {
			ret := multierror.Append(p.err, p.errs...)
			assert.NotNil(t, ret)
		})
	}
}

func TestAppendErr_Nil(t *testing.T) {
	for name, p := range tableNil {
		t.Run(name, func(t *testing.T) {
			ret := venom.AppendErr(p.err, p.errs...)
			assert.Nil(t, ret)
		})
	}
}

func TestAppendErr_NonNil(t *testing.T) {
	for name, p := range tableNonNil {
		t.Run(name, func(t *testing.T) {
			ret := venom.AppendErr(p.err, p.errs...)
			assert.NotNil(t, ret)
		})
	}
}
