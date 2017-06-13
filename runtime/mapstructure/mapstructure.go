package mapstructure

import (
	ms "github.com/mitchellh/mapstructure"
)

func Decode(input, output interface{}) error {
	config := &ms.DecoderConfig{
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: true,
		TagName:          "json",
	}

	decoder, err := ms.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}
