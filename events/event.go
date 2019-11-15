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

// type Event struct {
// 	Type    string      `json:"type" binding:"required" validate:"required"`
// 	Payload interface{} `json:"payload" binding:"required" validate:"required"`
// }

// func NewEvent(t string, p interface{}) *Event {
// 	return &Event{
// 		Type:    t,
// 		Payload: p,
// 	}
// }

// func FromBytes(b []byte) (*Event, errors.Error) {
// 	var e Event
// 	if err := json.Unmarshal(b, &e); err != nil {
// 		return nil, errors.NewInternal().SetPath("composition/events.FromBytes").SetCode("UNMARSHAL").SetMessage(err.Error())
// 	}

// 	return &e, nil

// }

// func (e *Event) ToBytes() ([]byte, errors.Error) {
// 	b, err := json.Marshal(e)
// 	if err != nil {
// 		return nil, errors.NewInternal().SetPath("composition/events.ToBytes").SetCode("MARSHAL").SetMessage(err.Error())
// 	}
// 	return b, nil
// }
