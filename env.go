package venom

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func parseEnvLine(line string) (string, string, error) {
	pair := strings.SplitN(line, "=", 2)
	if len(pair) != 2 {
		return "", "", fmt.Errorf("Invalid format, must be kv-pair separated with '=': %s", line)
	}
	return pair[0], pair[1], nil
}

func readEnvFile(p string) (map[string]string, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}

	s := bufio.NewScanner(f)
	envMap := make(map[string]string)
	for s.Scan() {
		line := s.Text()
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue // skip empty lines
		}

		k, v, err := parseEnvLine(line)
		if err != nil {
			return nil, err
		}
		envMap[k] = v
	}
	return envMap, nil
}

type envConfig struct {
	Envs      []string `mapstructure:"env"`
	EnvFiles  []string `mapstructure:"env-file"`
	EnvPrefix string   `mapstructure:"env-prefix"`
}

func initEnvFlags(flags *pflag.FlagSet, envprefix string) error {
	flags.String("env-prefix", envprefix, "Set environment variables prefix")
	flags.StringSliceP("env", "e", nil, "Set environment variables")
	flags.StringSlice("env-file", nil, "Read in a file of environment variables")
	return nil
}

func readEnv(v *viper.Viper) error {
	var cfg envConfig
	err := Unmarshal(&cfg, v)
	if err != nil {
		return err
	}

	envMap := make(map[string]string)

	// read --env-prefix
	if cfg.EnvPrefix != "" {
		v.SetEnvPrefix(cfg.EnvPrefix)
	}

	// read --env-file
	for _, f := range cfg.EnvFiles {
		env, err := readEnvFile(f)
		if err != nil {
			return err
		}
		for k, v := range env {
			envMap[k] = v
		}
	}

	// read --env
	for _, kv := range cfg.Envs {
		k, v, err := parseEnvLine(kv)
		if err != nil {
			return err
		}
		envMap[k] = v
	}

	// apply env
	for k, v := range envMap {
		err = os.Setenv(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func lookupEnvPrefix(flags *pflag.FlagSet) (string, bool) {
	envprefix := ""
	flags.VisitAll(func(f *pflag.Flag) {
		if f.Name == "env-prefix" {
			envprefix = f.Value.String()
		}
	})
	return envprefix, envprefix != ""
}
