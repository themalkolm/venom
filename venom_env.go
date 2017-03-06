package venom

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
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
	Envs     []string `mapstructure:"env"`
	EnvFiles []string `mapstructure:"env-file"`
}

func initEnvFlags(flags *pflag.FlagSet) error {
	var errors []string
	flags.VisitAll(func(f *pflag.Flag) {
		if f.Name == "env" || f.Name == "env-file" {
			errors = append(errors, fmt.Sprintf("Flag %s already exists!", f.Name))
		}
	})
	if len(errors) > 0 {
		return fmt.Errorf("%d errors:\n%s", len(errors), strings.Join(errors, "\n"))
	}

	flags.StringSliceP("env", "e", nil, "Set environment variables")
	flags.StringSlice("env-file", nil, "Read in a file of environment variables")
	return nil
}

//
// Viper does not decode string slices correctly
//
// https://github.com/spf13/viper/pull/319
//
func patchViper(stringSliceKeys []string, v *viper.Viper) {
	for _, k := range stringSliceKeys {
		if v.Get(k) == nil {
			continue // skip nil values
		}

		value := v.GetString(k)
		value = strings.TrimSpace(value)
		if value == "" {
			v.Set(k, []string{})
			continue // skip empty values
		}

		parts := strings.Split(value, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}
		v.Set(k, parts)
	}
}

func readEnv(viperMaybe ...*viper.Viper) error {
	v := viper.GetViper()
	if len(viperMaybe) != 0 {
		v = viperMaybe[0]
	}

	patchViper([]string{"env", "env-file"}, v)

	var cfg envConfig
	err := v.Unmarshal(&cfg)
	if err != nil {
		return err
	}

	envMap := make(map[string]string)

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

// Configure common flags and environment config considered (by me) as a good approach to _bootstrap_
// any 12-factor app.
//
// Note that we add some extra flags & alter PreRunE value.
//
func TwelveFactorCmd(cmd *cobra.Command, flags *pflag.FlagSet, viperMaybe ...*viper.Viper) error {
	v := viper.GetViper()
	if len(viperMaybe) != 0 {
		v = viperMaybe[0]
	}

	err := initEnvFlags(flags)
	if err != nil {
		return err
	}

	parts := strings.SplitN(cmd.Use, " ", 2)
	if len(parts) == 0 {
		return fmt.Errorf("Please set cmd.Use one-liner so name could be determined: %s", cmd.Use)
	}

	err = TwelveFactor(parts[0], flags, v)
	if err != nil {
		return err
	}

	// I have no idea if PreRunE is the right hook to use here.
	if cmd.PreRunE != nil {
		preRunE := cmd.PreRunE
		cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
			err := preRunE(cmd, args)
			if err != nil {
				return err
			}
			return readEnv(v)
		}
	} else {
		cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
			return readEnv(v)
		}
	}

	return nil
}
