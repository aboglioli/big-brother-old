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

// ArticleExistsEventRequest is a request to validate article
type ArticleExistsEventRequest struct {
	Type     string                           `json:"type"`
	Exchange string                           `json:"exchange"`
	Queue    string                           `json:"queue"`
	Message  articleExistsEventRequestMessage `json:"message"`
}

type articleExistsEventRequestMessage struct {
	ReferenceID string `json:"referenceId"`
	ArticleID   string `json:"articleId"`
}

// ArticleExistsEventResponse is a response from catalog with article validation
type ArticleExistsEventResponse struct {
	Type    string                           `json:"type"`
	Message articleExistsEventReponseMessage `json:"message"`
}

type articleExistsEventReponseMessage struct {
	ReferenceID string `json:"referenceId"`
	ArticleID   string `json:"articleId`
	Valid       bool   `json:"valid"`
}
