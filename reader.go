package venom

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

func readerFor(path string) (io.ReadCloser, error) {
	switch {
	case path == "-":
		return NopReadCloser(os.Stdin), nil
	default:
		return os.OpenFile(path, os.O_RDONLY, 0)
	}
}

func ReadObject(r io.Reader, format Format, out interface{}) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	switch format {
	case JSONFormat:
		return json.Unmarshal(b, out)
	case YAMLFormat:
		return yaml.Unmarshal(b, out)
	default:
		return fmt.Errorf("Unrecognized format: %s", format)
	}
}

func ReadObjectFrom(path string, out interface{}) error {
	r, err := readerFor(path)
	if err != nil {
		return err
	}
	defer r.Close()

	switch {
	case path == "-":
		return ReadObject(r, DefaultInputFormat, out)
	case strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml"):
		return ReadObject(r, YAMLFormat, out)
	case strings.HasSuffix(path, ".json"):
		return ReadObject(r, JSONFormat, out)
	default:
		return fmt.Errorf("Can't deduce file format: %s", path)
	}
}
