package venom

import (
	"fmt"
	"strings"
)

func parseMapStringString(s, sep, kvsep string) (map[string]string, error) {
	if s == "" {
		return make(map[string]string, 0), nil
	}

	parts := strings.Split(s, sep)
	m := make(map[string]string, len(parts))
	for _, s := range parts {
		pair := strings.SplitN(s, kvsep, 2)
		if len(pair) != 2 {
			return nil, fmt.Errorf("Invalid kv pair, must be seprated by %s but got: %s", kvsep, s)
		}
		k := pair[0]
		v := pair[1]
		m[k] = v
	}
	return m, nil
}
