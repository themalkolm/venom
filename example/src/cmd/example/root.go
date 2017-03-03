package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"github.com/themalkolm/venom"
)

type Config struct {
	Foo    string `mapstructure:"foo"`
	FooBar string `mapstructure:"foo-bar"`
}

var RootCmd = &cobra.Command{
	Use:          "example",
	Short:        "Do example things.",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg Config
		err := viper.Unmarshal(&cfg)
		if err != nil {
			return err
		}

		return runE(&cfg)
	},
}

func init() {
	RootCmd.PersistentFlags().String("foo", "", "Some foonees must be set")
	RootCmd.PersistentFlags().String("foo-bar", "", "Some barness must be set")

	err := venom.TwelveFactor("example", RootCmd.PersistentFlags())
	if err != nil {
		log.Fatal(err)
	}
}

func runE(cfg *Config) error {
	b, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}
	fmt.Printf("%+v", string(b))
	return nil
}
