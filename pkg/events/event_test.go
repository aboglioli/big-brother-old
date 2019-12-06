package events

import (
	"testing"

	"github.com/aboglioli/big-brother/pkg/converter"
	"github.com/aboglioli/big-brother/pkg/tests/assert"
)

type customEvent struct {
	Event
	Data string
}

func TestCustomEventEncoding(t *testing.T) {
	cEvt := &customEvent{Event{"Custom"}, "This is data"}
	conv := converter.DefaultConverter()

	src, err := conv.Encode(cEvt)
	assert.Ok(t, err)

	var dst customEvent
	assert.Ok(t, conv.Decode(src, &dst))

	assert.Assert(t, dst.Type == cEvt.Type && dst.Data == cEvt.Data)
}
