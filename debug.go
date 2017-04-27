package venom

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

//
// Debug configuration for every venom command. Makes it easier to debug and figure
// out how exactly all env, flags etc. are merged.
//
type debugConfig struct {
	PrintConfig bool `mapstructure:"print-config"`
}

func initDebugFlags(flags *pflag.FlagSet) error {
	var errors []string
	flags.VisitAll(func(f *pflag.Flag) {
		if f.Name == "print-config" {
			errors = append(errors, fmt.Sprintf("Flag %s already defined!", f.Name))
		}
	})
	if len(errors) > 0 {
		return fmt.Errorf("%d errors:\n%s", len(errors), strings.Join(errors, "\n"))
	}

	flags.Bool("print-config", false, "Print result configuraiton and exit.")
	return nil
}

func readDebug(v *viper.Viper) error {
	var cfg debugConfig
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
