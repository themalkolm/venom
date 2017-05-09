package venom

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	//
	// Hide flags that are used to generate result configuration itself.
	//
	hideFlags = []string{"print-config", "print-env", "env", "env-file"}
)

//
// Debug configuration for every venom command. Makes it easier to debug and figure
// out how exactly all env, flags etc. are merged.
//
type debugConfig struct {
	PrintConfig bool `mapstructure:"print-config"`
	PrintEnv    bool `mapstructure:"print-env"`
}

func initDebugFlags(flags *pflag.FlagSet) error {
	flags.Bool("print-config", false, "Print result configuraiton and exit.")
	flags.Bool("print-env", false, "Print env and exit.")
	return nil
}

func readDebug(flags *pflag.FlagSet, v *viper.Viper) error {
	var cfg debugConfig
	err := Unmarshal(&cfg, v)
	if err != nil {
		return err
	}

	if cfg.PrintConfig {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "    ")

		all := v.AllSettings()
		for _, k := range hideFlags {
			delete(all, k)
		}

		err := enc.Encode(all)
		if err != nil {
			return err
		}

		os.Exit(0)
	}

	if cfg.PrintEnv {
		all := v.AllSettings()
		for _, k := range hideFlags {
			delete(all, k)
		}

		keys := make([]string, 0, len(all))
		for k := range all {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		envprefix, found := lookupEnvPrefix(flags)
		if !found {
			envprefix = "?"
		}

		r := envKeyReplacer(flags)
		for _, k := range keys {
			v := v.Get(k)
			fmt.Printf("%s_%s=%v\n", envprefix, r.Replace(strings.ToUpper(k)), v)
		}

		os.Exit(0)
	}

	return nil
}
