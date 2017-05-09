package venom

import (
	"sort"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func allKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

//
// Better version of viper.AutomaticEnv that searches FOO_BAR for every --foo-bar key in
// addition to the default FOO-BAR.
//
// Note that it must be called *after* all flags are added.
//
func automaticEnv(flags *pflag.FlagSet, v *viper.Viper) {
	replaceMap := make(map[string]string, flags.NFlag())
	flags.VisitAll(func(f *pflag.Flag) {
		name := strings.ToUpper(f.Name)
		replaceMap[name] = sanitize(name)
	})

	keys := allKeys(replaceMap)

	// Reverse sort keys, this is to make sure foo-bar comes before foo. This is to prevent
	// foo being triggered when foo-bar is given to string replacer.
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	values := make([]string, 0, 2*len(keys))
	for _, k := range keys {
		values = append(values, k, replaceMap[k])
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(values...))
	v.AutomaticEnv()
}
