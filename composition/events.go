package composition

// CompositionChangedEvent is a single composition change
type CompositionChangedEvent struct {
	Type        string       `json:"type"`
	Composition *Composition `json:"composition"`
}

// CompositionUpdatedAutomatically is an event containing all updated compositions after a dependency update
type CompositionsUpdatedAutomaticallyEvent struct {
	Type         string         `json:"type"`
	Compositions []*Composition `json:"compositions"`
}
