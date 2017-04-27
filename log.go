package venom

import (
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	LogLevels  = []string{}
	LogFormats = []string{"text", "json"}
)

func init() {
	for _, l := range logrus.AllLevels {
		LogLevels = append(LogLevels, l.String())
	}
}

type logConfig struct {
	LogLevel     string `mapstructure:"log-level"`
	LogFormatter string `mapstructure:"log-format"`
}

func initLogFlags(flags *pflag.FlagSet) error {
	var errors []string
	flags.VisitAll(func(f *pflag.Flag) {
		if f.Name == "log-level" || f.Name == "log-format" {
			errors = append(errors, fmt.Sprintf("Flag %s already defined!", f.Name))
		}
	})
	if len(errors) > 0 {
		return fmt.Errorf("%d errors:\n%s", len(errors), strings.Join(errors, "\n"))
	}

	flags.String("log-level", "info", fmt.Sprintf("Log level [%s]", strings.Join(LogLevels, "|")))
	flags.String("log-format", "text", fmt.Sprintf("Log format [%s]", strings.Join(LogFormats, "|")))
	return nil
}

func readLog(viperMaybe ...*viper.Viper) error {
	v := viper.GetViper()
	if len(viperMaybe) != 0 {
		v = viperMaybe[0]
	}

	var cfg logConfig
	err := v.Unmarshal(&cfg)
	if err != nil {
		return err
	}

	l, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		return err
	}

	switch cfg.LogFormatter {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{})
	default:
		return fmt.Errorf("Invalid log format: %s", cfg.LogFormatter)
	}

	logrus.SetLevel(l)
	return nil
}
