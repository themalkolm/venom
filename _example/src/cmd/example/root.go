package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

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
		err := viper.Unmarshal(&cfg)
		if err != nil {
			return err
		}
		return runE(&cfg)
	},
}

type Config struct {
	Foo    string `mapstructure:"foo"`
	FooBar string `mapstructure:"foo-bar"`
	FooMoo int    `mapstructure:"foo-moo" pflag:"foo-moo,m,Some mooness must be set"`
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
	flags := venom.DefineFlags(Config{})
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
	// with 12-factor, just make it easy to play with.
	//
	b, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}
	fmt.Printf("%+v", string(b))
	return nil
}
