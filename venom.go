package venom

import (
	"sort"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func sanitize(s string) string {
	return strings.Replace(s, "-", "_", -1)
}

//
// Better version of viper.AutomaticEnv that searches FOO_BAR for every --foo-bar key in
// addition to the default FOO-BAR.
//
func AutomaticEnv(flags *pflag.FlagSet, viperMaybe ...*viper.Viper) {
	v := viper.GetViper()
	if len(viperMaybe) != 0 {
		v = viperMaybe[0]
	}

	replaceMap := make(map[string]string, flags.NFlag())
	flags.VisitAll(func(f *pflag.Flag) {
		name := strings.ToUpper(f.Name)
		replaceMap[name] = sanitize(name)
	})

	keys := make([]string, 0, len(replaceMap))
	for k := range replaceMap {
		keys = append(keys, k)
	}

	// Reverse sort keys, this is to make sure foo-bar comes before foo. This is to prevent
	// foo being triggered when foo-bar is given to string replacer.
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	values := make([]string, 0, 2*len(keys))
	for _, k := range keys {
		values = append(values, k)
		values = append(values, replaceMap[k])
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(values...))
	v.AutomaticEnv()
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
	AutomaticEnv(flags, v)

	return nil
}
