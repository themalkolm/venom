package venom

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func preRun(flags *pflag.FlagSet, v *viper.Viper) error {
	err := readEnv(v)
	if err != nil {
		return err
	}

	err = readLog(v)
	if err != nil {
		return err
	}

	// must be the last one
	err = readDebug(flags, v)
	if err != nil {
		return err
	}

	return nil
}

//
// Configure common flags and environment config considered (by me) as a good approach to _bootstrap_
// any 12-factor app.
//
// Note that we add some extra flags & alter PreRunE value.
//
func TwelveFactorCmd(name string, cmd *cobra.Command, flags *pflag.FlagSet, viperMaybe ...*viper.Viper) error {
	v := viper.GetViper()
	if len(viperMaybe) != 0 {
		v = viperMaybe[0]
	}

	if name == "" {
		parts := strings.SplitN(cmd.Use, " ", 2)
		if len(parts) == 0 {
			return fmt.Errorf("Please either provide name or set cmd.Use one-liner so name could be determined: %s", cmd.Use)
		}
		name = parts[0]
	}

	err := TwelveFactor(name, flags, v)
	if err != nil {
		return err
	}

	// I have no idea if PersistentPreRunE is the right hook to use here.
	if cmd.PersistentPreRunE != nil {
		preRunE := cmd.PersistentPreRunE
		cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
			err := preRun(flags, v)
			if err != nil {
				return err
			}

			return preRunE(cmd, args)
		}
	} else {
		cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
			return preRun(flags, v)
		}
	}

	go ListenAndServeMetrics(name, v)

	if cmd.RunE != nil {
		cmd.RunE = CronRunE(cmd.RunE, v)
	}

	return nil
}

//
// Version of TwelveFactorCmd that parses flags from the provided defaults.
//
func AutoFlagsTwelveFactorCmd(name string, cmd *cobra.Command, defaults interface{}, viperMaybe ...*viper.Viper) error {
	flags := cmd.PersistentFlags()
	flags.AddFlagSet(MustDefineFlags(defaults))

	return TwelveFactorCmd(name, cmd, flags, viperMaybe...)
}
