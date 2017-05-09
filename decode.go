package venom

import (
	"github.com/spf13/viper"
)

//
// Config structure that is able to self-validate.
//
type Config interface {
	Valid() error
}

//
// Decode config from viper and self-validate.
//
func Decode(out Config, cfg *viper.Viper) error {
	err := Unmarshal(out, cfg)
	if err != nil {
		return err
	}
	return out.Valid()
}
