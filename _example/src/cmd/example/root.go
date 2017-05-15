package main

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"github.com/Sirupsen/logrus"
	"github.com/themalkolm/venom"
)

var RootCmd = &cobra.Command{
	Use:          "example",
	Short:        "Do example things.",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		//
		// 12-factor automatically registered all flags in viper configuration. This
		// means viper is already configured with all values provided via cli, env
		// and defaults if any: flag > env > default.
		//
		// The following trick with having a struct representing our configuration
		// allows us to escape dynamic nature of viper and have a static Config
		// structure.
		//
		var cfg Config
		err := venom.Unmarshal(&cfg, viper.GetViper())
		if err != nil {
			return err
		}
		return runE(&cfg)
	},
}

type Inner struct {
	Goo string `mapstructure:"goo" pflag:"goo,,Some gooness must be set"`
}

type Config struct {
	Inner `mapstructure:",squash" pflag:"++"`

	Foo        string            `mapstructure:"foo"`
	FooBar     string            `mapstructure:"foo-bar"`
	FooMoo     int               `mapstructure:"foo-moo"  pflag:"foo-moo,m,Some mooness must be set"`
	Time       time.Time         `mapstructure:"time" pflag:"time,,Some time"`
	Duration   time.Duration     `mapstructure:"duration" pflag:"duration,,Some duration"`
	Bools      []bool            `mapstructure:"bools"  pflag:"bools,,Some bools"`
	Ints       []int             `mapstructure:"ints"  pflag:"ints,,Some ints"`
	Uints      []uint            `mapstructure:"uints"  pflag:"uints,,Some uints"`
	Strings    []string          `mapstructure:"strings"  pflag:"strings,,Some strings"`
	StringsMap map[string]string `mapstructure:"strings-map"  pflag:"strings-map,,Some strings map"`
}

func init() {
	//
	// You either define flags manually ...
	//
	RootCmd.PersistentFlags().String("foo", "", "Some foonees must be set")
	RootCmd.PersistentFlags().String("foo-bar", "", "Some barness must be set")

	//
	// ... or let venom to do it for you.
	//
	defaults := Config{
		FooMoo: 43,
	}
	flags := venom.MustDefineFlags(defaults)
	RootCmd.PersistentFlags().AddFlagSet(flags)

	//
	// Enable 12-factor application so magic happens
	//
	err := venom.TwelveFactorCmd("example", RootCmd, RootCmd.PersistentFlags())
	if err != nil {
		log.Fatal(err)
	}
}

func runE(cfg *Config) error {
	//
	// Here we simply dump passed config object. This has nothing to do
	// with 12-factor, just to make it easy to play with.
	//
	b, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}
	fmt.Printf("%+v", string(b))

	logrus.Info("Logging using [info] level")
	logrus.Warn("Logging using [warning] level")
	logrus.Debug("Logging using [debug] level")
	logrus.Error("Logging using [error] level")
	logrus.Fatal("Logging using [fatal] level")
	logrus.Panic("Logging using [panic] level")

	return nil
}
