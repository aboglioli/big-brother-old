package converter

import (
	"encoding/json"

	"github.com/aboglioli/big-brother/pkg/errors"
)

// Converter is an interface
type Converter interface {
	Decode(src []byte, dst interface{}) error
	Encode(src interface{}) ([]byte, error)
}

// Default converter for json structures
type jsonConverter struct {
}

func DefaultConverter() *jsonConverter {
	return &jsonConverter{}
}

func (c *jsonConverter) Decode(src []byte, dst interface{}) error {
	if err := json.Unmarshal(src, dst); err != nil {
		return errors.NewInternal("FAILED_TO_DECODE").SetPath("infrastructure/events/rabbit.Decode").SetRef(err)
	}

	return nil
}
func (c *jsonConverter) Encode(src interface{}) ([]byte, error) {
	b, err := json.Marshal(src)
	if err != nil {
		return nil, errors.NewInternal("FAILTED_TO_CODE").SetPath("infrastructure/events/rabbit.Code").SetRef(err)
	}
	return b, nil
}
