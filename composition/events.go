package composition

import (
	"github.com/aboglioli/big-brother/events"
)

// CompositionChangedEvent is a single composition change
type CompositionChangedEvent struct {
	events.Event
	Composition *Composition `json:"composition"`
}

// CompositionUpdatedAutomatically is an event containing all updated compositions after a dependency update
type CompositionsUpdatedAutomaticallyEvent struct {
	events.Event
	Compositions []*Composition `json:"compositions"`
}
