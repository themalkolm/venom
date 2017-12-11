package venom

import (
	"time"

	"github.com/mitchellh/mapstructure"
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

//
// Same as viper.Unmarshal but with support for string slices, time etc.
//
func Unmarshal(out interface{}, cfg *viper.Viper) error {
	return decode(cfg.AllSettings(), defaultDecoderConfig(out))
}

func defaultDecoderConfig(output interface{}) *mapstructure.DecoderConfig {
	return &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			stringToStringSliceHookFunc(","),
			stringToBoolSliceHookFunc(","),
			stringToIntSliceHookFunc(","),
			stringToUintSliceHookFunc(","),
			stringToTimeDurationHookFunc(),
			stringToTimeHookFunc(time.RFC3339),
			stringToMapStringStringHookFunc(",", "="),
		),
	}
}

func decode(input interface{}, config *mapstructure.DecoderConfig) error {
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}
