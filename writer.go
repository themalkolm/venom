package venom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"
	"strings"
)

type OutputFormat string

var (
	OutputJSONFormat = OutputFormat("json")
	OutputYAMLFormat = OutputFormat("yaml")
	OutputFormats    = []OutputFormat{
		OutputJSONFormat,
		OutputYAMLFormat,
	}

	DefaultOutputFormat = OutputYAMLFormat
)

func writerFor(path string) (io.WriteCloser, error) {
	switch {
	case path == "":
		return NopWriteCloser(os.Stdout), nil
	default:
		return os.OpenFile(path, os.O_WRONLY, 0644)
	}
}

func WriteObject(in interface{}, format OutputFormat, w io.Writer) error {
	switch format {
	case OutputYAMLFormat:
		b, err := yaml.Marshal(in)
		if err != nil {
			return err
		}
		_, err = w.Write(b)
		return err
	case OutputJSONFormat:
		var b bytes.Buffer

		e := json.NewEncoder(&b)
		e.SetIndent("", "    ")

		err := e.Encode(in)
		if err != nil {
			return err
		}
		_, err = io.Copy(w, &b)
		return err
	default:
		return fmt.Errorf("Unrecognized format: %s", format)
	}
}

func WriteObjectTo(in interface{}, path string) error {
	w, err := writerFor(path)
	if err != nil {
		return err
	}
	defer w.Close()

	switch {
	case path == "-":
		return WriteObject(in, DefaultOutputFormat, w)
	case strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml"):
		return WriteObject(in, OutputYAMLFormat, w)
	case strings.HasSuffix(path, ".json"):
		return WriteObject(in, OutputJSONFormat, w)
	default:
		return fmt.Errorf("Can't deduce file format: %s", path)
	}
}
