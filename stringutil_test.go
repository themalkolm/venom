package venom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type setup struct {
	s string
	m map[string]string
}

var (
	table = map[string]setup{
		"empty": {
			s: "",
			m: map[string]string{},
		},
		"single": {
			s: "foo=bar",
			m: map[string]string{
				"foo": "bar",
			},
		},
		"multiple": {
			s: "foo=bar,goo=moo",
			m: map[string]string{
				"foo": "bar",
				"goo": "moo",
			},
		},
	}
)

func TestParse(t *testing.T) {
	for name, x := range table {
		t.Run(name, func(t *testing.T) {
			m, err := parseMapStringString(x.s, ",", "=")
			assert.Nil(t, err)
			assert.Equal(t, m, x.m)
		})
	}
}

func TestSerialize(t *testing.T) {
	for name, x := range table {
		t.Run(name, func(t *testing.T) {
			s, err := serializeMapStringString(x.m, ",", "=")
			assert.Nil(t, err)
			assert.Equal(t, s, x.s)
		})
	}
}

func TestFailNoSep(t *testing.T) {
	m, err := parseMapStringString("foo", ",", "=")
	assert.NotNil(t, err)
	assert.Nil(t, m)
}
