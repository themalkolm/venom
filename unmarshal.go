package venom

import (
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

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
