package venom

import (
	"github.com/spf13/viper"
)

type Config interface {
	Valid() error
}

func Decode(out Config, cfg *viper.Viper) error {
	err := Unmarshal(out, cfg)
	if err != nil {
		return err
	}
	return out.Valid()
}
