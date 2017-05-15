package venom

import (
	"bytes"
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

func serializeMapStringString(m map[string]string, sep, kvsep string) (string, error) {
	records := make([]string, 0, len(m))
	for k, v := range m {
		records = append(records, fmt.Sprintf("%s%s%s", k, kvsep, v))
	}

	b := &bytes.Buffer{}
	w := csv.NewWriter(b)
	w.Comma = []rune(sep)[0]
	err := w.Write(records)
	if err != nil {
		return "", err
	}
	w.Flush()
	return strings.TrimSuffix(b.String(), "\n"), nil
}
