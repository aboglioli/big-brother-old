package events

import (
	"testing"

	"github.com/aboglioli/big-brother/pkg/tests"
)

type customEvent struct {
	Event
	Data string
}

func TestCustomEventEncoding(t *testing.T) {
	cEvt := &customEvent{Event{"Custom"}, "This is data"}
	conv := DefaultConverter()

	src, err := conv.Encode(cEvt)
	tests.Ok(t, err)

	var dst customEvent
	tests.Ok(t, conv.Decode(src, &dst))

	tests.Assert(t, dst.Type == cEvt.Type && dst.Data == cEvt.Data)
}
