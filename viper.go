package venom

import (
	"sort"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

//
// Better version of viper.AutomaticEnv that searches FOO_BAR for every --foo-bar key in
// addition to the default FOO-BAR.
//
func automaticEnv(flags *pflag.FlagSet, v *viper.Viper) {
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
// Viper does not decode string slices correctly
//
// https://github.com/spf13/viper/pull/319
//
func parseStringSliceFlags(flagNames []string, v *viper.Viper) {
	for _, k := range flagNames {
		if v.Get(k) == nil {
			continue // skip nil values
		}

		value := v.GetString(k)
		parts := strings.Split(value, ",")

		values := []string{}
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" { // skip empty values
				values = append(values, p)
			}
		}

		v.Set(k, values)
	}
}
