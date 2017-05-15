package venom

import (
	"encoding/csv"
	"fmt"
	"strings"
)

func parseMapStringString(s, sep, kvsep string) (map[string]string, error) {
	if s == "" {
		return make(map[string]string, 0), nil
	}

	r := csv.NewReader(strings.NewReader(s))
	r.Comma = []rune(sep)[0]
	records, err := r.Read()
	if err != nil {
		return nil, err
	}

	m := make(map[string]string, len(records))
	for _, record := range records {
		pair := strings.SplitN(record, kvsep, 2)
		if len(pair) != 2 {
			return nil, fmt.Errorf("Invalid kv pair, must be seprated by '%s' but got: %s", kvsep, s)
		}
		k := pair[0]
		v := pair[1]
		m[k] = v
	}
	return m, nil
}
