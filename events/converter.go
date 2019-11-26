package events

import (
	"encoding/json"

	"github.com/aboglioli/big-brother/errors"
)

// Default converter
type jsonConverter struct {
}

func DefaultConverter() *jsonConverter {
	return &jsonConverter{}
}

func (c *jsonConverter) Decode(src []byte, dst interface{}) errors.Error {
	if err := json.Unmarshal(src, dst); err != nil {
		return errors.NewInternal().SetPath("infrastructure/events/rabbit.Decode").SetCode("FAILED_TO_DECODE").SetMessage(err.Error())
	}

	return nil
}
func (c *jsonConverter) Code(src interface{}) ([]byte, errors.Error) {
	b, err := json.Marshal(src)
	if err != nil {
		return nil, errors.NewInternal().SetPath("infrastructure/events/rabbit.Code").SetCode("FAILTED_TO_CODE").SetMessage(err.Error())
	}
	return b, nil
}
