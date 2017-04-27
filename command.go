package venom

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

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

func preRun(viperMaybe ...*viper.Viper) error {
	v := viper.GetViper()
	if len(viperMaybe) != 0 {
		v = viperMaybe[0]
	}

	var cfg commonConfig
	err := v.Unmarshal(&cfg)
	if err != nil {
		return err
	}

	if cfg.PrintConfig {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "    ")

		all := v.AllSettings()
		for _, k := range []string{"print-config", "env", "env-file"} {
			delete(all, k)
		}

		err := enc.Encode(all)
		if err != nil {
			return err
		}

		os.Exit(0)
	}

	return nil
}

//
// Configure common flags and environment config considered (by me) as a good approach to _bootstrap_
// any 12-factor app.
//
// Note that we add some extra flags & alter PreRunE value.
//
func TwelveFactorCmd(name string, cmd *cobra.Command, flags *pflag.FlagSet, viperMaybe ...*viper.Viper) error {
	v := viper.GetViper()
	if len(viperMaybe) != 0 {
		v = viperMaybe[0]
	}

	err := initCommonFlags(flags)
	if err != nil {
		return err
	}

	err = initEnvFlags(flags)
	if err != nil {
		return err
	}

	err = initLogFlags(flags)
	if err != nil {
		return err
	}

	if name == "" {
		parts := strings.SplitN(cmd.Use, " ", 2)
		if len(parts) == 0 {
			return fmt.Errorf("Please either provide name or set cmd.Use one-liner so name could be determined: %s", cmd.Use)
		}
		name = parts[0]
	}

	err = TwelveFactor(name, flags, v)
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

			err = readEnv(v)
			if err != nil {
				return err
			}

			err = readLog(v)
			if err != nil {
				return err
			}

			return preRun(v)
		}
	} else {
		cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
			err := readEnv(v)
			if err != nil {
				return err
			}
			err = readLog(v)
			if err != nil {
				return err
			}
			return preRun(v)
		}
	}

	return nil
}
