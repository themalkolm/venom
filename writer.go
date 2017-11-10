package venom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Format string

var (
	JSONFormat = Format("json")
	YAMLFormat = Format("yaml")

	InputFormats = []Format{
		JSONFormat,
		YAMLFormat,
	}
	OutputFormats = []Format{
		JSONFormat,
		YAMLFormat,
	}

	DefaultInputFormat  = YAMLFormat
	DefaultOutputFormat = YAMLFormat
)

func writerFor(path string) (io.WriteCloser, error) {
	switch {
	case path == "":
		return NopWriteCloser(os.Stdout), nil
	default:
		return os.OpenFile(path, os.O_WRONLY, 0644)
	}
}

func WriteObject(in interface{}, format Format, w io.Writer) error {
	switch format {
	case YAMLFormat:
		b, err := yaml.Marshal(in)
		if err != nil {
			return err
		}
		_, err = w.Write(b)
		return err
	case JSONFormat:
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
	case path == "":
		return WriteObject(in, DefaultOutputFormat, w)
	case strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml"):
		return WriteObject(in, YAMLFormat, w)
	case strings.HasSuffix(path, ".json"):
		return WriteObject(in, JSONFormat, w)
	default:
		return fmt.Errorf("Can't deduce file format: %s", path)
	}
}

//
// Simple wrapper around any object that could be used in logrus e.g.
//
// logrus.WithField("plan", venom.PrettyField(plan)).Info("Here is the plan")
//
func PrettyField(in interface{}, format Format) string {
	var b bytes.Buffer
	err := WriteObject(in, format, &b)
	if err != nil {
		return fmt.Sprintf("failed to pretty print: %s", err)
	}
	return b.String()
}
