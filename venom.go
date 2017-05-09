package venom

import (
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

//
// 12-factor setup for viper-backed application.
//
// Most important it will check environment variables for all flags in the provided flag set.
//
func TwelveFactor(name string, flags *pflag.FlagSet, viperMaybe ...*viper.Viper) error {
	v := viper.GetViper()
	if len(viperMaybe) != 0 {
		v = viperMaybe[0]
	}

	err := initDebugFlags(flags)
	if err != nil {
		return err
	}

	envprefix := strings.ToUpper(sanitize(name))
	err = initEnvFlags(flags, envprefix)
	if err != nil {
		return err
	}

	err = initLogFlags(flags)
	if err != nil {
		return err
	}

	// Bind flags and configuration keys 1-to-1
	err = v.BindPFlags(flags)
	if err != nil {
		return err
	}

	// Set env prefix
	v.SetEnvPrefix(envprefix)

	// Patch automatic env
	AutomaticEnv(flags, v)

	return nil
}
