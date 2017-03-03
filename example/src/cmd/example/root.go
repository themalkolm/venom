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
		b, err := yaml.Marshal(viper.AllSettings())
		if err != nil {
			return err
		}
		fmt.Printf("%+v", string(b))
		return nil
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
