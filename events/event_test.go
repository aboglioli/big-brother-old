package events

import (
	"testing"
)

type customEvent struct {
	Event
	Data string
}

func TestCustomEventEncoding(t *testing.T) {
	cEvt := &customEvent{Event{"Custom"}, "This is data"}
	conv := DefaultConverter()

	src, err := conv.Encode(cEvt)
	if err != nil {
		t.Error(err)
	}

	var dst customEvent
	if err := conv.Decode(src, &dst); err != nil {
		t.Error(err)
	}

	if dst.Type != cEvt.Type || dst.Data != cEvt.Data {
		t.Error("Decoded and encoded events are different")
	}
}
