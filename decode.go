package venom

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func defaultDecoderConfig(output interface{}) *mapstructure.DecoderConfig {
	return &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: true,
		DecodeHook:       mapstructure.StringToTimeDurationHookFunc(),
	}
}

func decode(input interface{}, config *mapstructure.DecoderConfig) error {
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}

func Unmarshal(rawVal interface{}, v *viper.Viper) error {
	return decode((*viper.Viper)(v).AllSettings(), defaultDecoderConfig(rawVal))
}
