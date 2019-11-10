package composition

import (
	"encoding/json"

	"github.com/aboglioli/big-brother/errors"
)

type Event struct {
	Type        string       `json:"type"`
	Composition *Composition `json:"composition"`
}

func NewEvent(t string, c *Composition) *Event {
	return &Event{
		Type:        t,
		Composition: c,
	}
}

func FromBytes(b []byte) (*Event, errors.Error) {
	var e Event
	if err := json.Unmarshal(b, &e); err != nil {
		return nil, errors.New("composition/events.FromBytes", "UNMARSHAL", err.Error())
	}

	return &e, nil

}

func (e *Event) ToBytes() ([]byte, errors.Error) {
	b, err := json.Marshal(e)
	if err != nil {
		return nil, errors.New("composition/events.ToBytes", "MARSHAL", err.Error())
	}
	return b, nil
}
