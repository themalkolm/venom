package venom

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
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
