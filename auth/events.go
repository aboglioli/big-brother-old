package auth

import (
	"encoding/json"

	"github.com/aboglioli/big-brother/errors"
)

type LogoutEvent struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (e *LogoutEvent) FromBytes(b []byte) errors.Error {
	if err := json.Unmarshal(b, e); err != nil {
		return errors.NewInternal().SetPath("composition/events.FromBytes").SetCode("UNMARSHAL").SetMessage(err.Error())
	}

	return nil
}

func (e *LogoutEvent) ToBytes() ([]byte, errors.Error) {
	b, err := json.Marshal(e)
	if err != nil {
		return nil, errors.NewInternal().SetPath("composition/events.ToBytes").SetCode("MARSHAL").SetMessage(err.Error())
	}
	return b, nil
}
