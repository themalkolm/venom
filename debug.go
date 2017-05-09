package venom

import (
	"encoding/json"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	//
	// Hide flags that are used to generate result configuration itself.
	//
	hideFlags = []string{"print-config", "env", "env-file"}
)

//
// Debug configuration for every venom command. Makes it easier to debug and figure
// out how exactly all env, flags etc. are merged.
//
type debugConfig struct {
	PrintConfig bool `mapstructure:"print-config"`
}

func initDebugFlags(flags *pflag.FlagSet) error {
	flags.Bool("print-config", false, "Print result configuraiton and exit.")
	return nil
}

func readDebug(v *viper.Viper) error {
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

	return nil
}
