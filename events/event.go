package events

import (
	"encoding/json"

	"github.com/aboglioli/big-brother/errors"
)

type Type struct {
	Type string `json:"type"`
}

func (e *Type) FromBytes(b []byte) errors.Error {
	if err := json.Unmarshal(b, e); err != nil {
		return errors.NewInternal().SetPath("events/event.FromBytes").SetCode("UNMARSHAL").SetMessage(err.Error())
	}

	return nil
}

type Converter interface {
	FromBytes(b []byte) errors.Error
	ToBytes() ([]byte, errors.Error)
}
