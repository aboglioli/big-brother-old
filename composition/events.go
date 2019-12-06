package composition

import (
	"github.com/aboglioli/big-brother/pkg/events"
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

func NewCompositionCreatedEvent(c *Composition) (*CompositionChangedEvent, *events.Options) {
	event := &CompositionChangedEvent{events.Event{"CompositionCreated"}, c}
	opts := &events.Options{"composition", "topic", "composition.created", ""}
	return event, opts
}

func NewCompositionUpdatedManuallyEvent(c *Composition) (*CompositionChangedEvent, *events.Options) {
	event := &CompositionChangedEvent{events.Event{"CompositionUpdatedManually"}, c}
	opts := &events.Options{"composition", "topic", "composition.updated", ""}
	return event, opts
}

func NewCompositionDeletedEvent(c *Composition) (*CompositionChangedEvent, *events.Options) {
	event := &CompositionChangedEvent{events.Event{"CompositionDeleted"}, c}
	opts := &events.Options{"composition", "topic", "composition.deleted", ""}
	return event, opts
}

func NewCompositionsUpdatedAutomaticallyEvent(comps []*Composition) (*CompositionsUpdatedAutomaticallyEvent, *events.Options) {
	event := &CompositionsUpdatedAutomaticallyEvent{events.Event{"CompositionsUpdatedAutomatically"}, comps}
	opts := &events.Options{"composition", "topic", "composition.updated", ""}
	return event, opts
}

func NewCompositionUsesUpdatedSinceLastChangeEvent(c *Composition) (*CompositionChangedEvent, *events.Options) {
	event := &CompositionChangedEvent{events.Event{"CompositionUsesUpdatedSinceLastChange"}, c}
	opts := &events.Options{"composition", "topic", "composition.updated", ""}
	return event, opts
}
