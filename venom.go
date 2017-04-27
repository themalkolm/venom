package venom

import (
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func sanitize(s string) string {
	return strings.Replace(s, "-", "_", -1)
}

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

	// Bind flags and configuration keys 1-to-1
	err := v.BindPFlags(flags)
	if err != nil {
		return err
	}

	// Set env prefix
	v.SetEnvPrefix(strings.ToUpper(sanitize(name)))

	// Patch automatic env
	automaticEnv(flags, v)

	return nil
}
