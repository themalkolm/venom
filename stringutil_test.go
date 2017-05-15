package venom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEmpty(t *testing.T) {
	m, err := parseMapStringString("", ",", "=")
	assert.Nil(t, err)
	assert.Equal(t, m, map[string]string{})
}

func TestParseSingle(t *testing.T) {
	m, err := parseMapStringString("foo=bar", ",", "=")
	assert.Nil(t, err)
	assert.Equal(t, m, map[string]string{
		"foo": "bar",
	})
}

func TestParseMultiple(t *testing.T) {
	m, err := parseMapStringString("foo=bar,moo=goo", ",", "=")
	assert.Nil(t, err)
	assert.Equal(t, m, map[string]string{
		"foo": "bar",
		"moo": "goo",
	})
}

func TestFailNoSep(t *testing.T) {
	m, err := parseMapStringString("foo", ",", "=")
	assert.NotNil(t, err)
	assert.Nil(t, m)
}
