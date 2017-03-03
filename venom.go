package venom

import (
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Teach viper to search FOO_BAR for every --foo-bar key instead of
// the default FOO-BAR.
func AutomaticEnv(flags *pflag.FlagSet, viperMaybe ...*viper.Viper) {
	v := viper.GetViper()
	if len(viperMaybe) != 0 {
		v = viperMaybe[0]
	}

	replaceMap := make([]string, 0)
	flags.VisitAll(func(f *pflag.Flag) {
		name := strings.ToUpper(f.Name)
		replaceMap = append(replaceMap, name)
		replaceMap = append(replaceMap, strings.Replace(name, "-", "_", -1))
	})
	v.SetEnvKeyReplacer(strings.NewReplacer(replaceMap...))
	v.AutomaticEnv()
}

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
	v.SetEnvPrefix(strings.ToUpper(name))

	// Patch automatic env
	AutomaticEnv(flags, v)

	return nil
}
