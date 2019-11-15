package composition

import (
	"encoding/json"

	"github.com/aboglioli/big-brother/errors"
)

type Type struct {
	Type string `json:"type"`
}

func (e *Type) FromBytes(b []byte) errors.Error {
	if err := json.Unmarshal(b, e); err != nil {
		return errors.NewInternal().SetPath("composition/events.FromBytes").SetCode("UNMARSHAL").SetMessage(err.Error())
	}

	return nil
}

// CompositionChangedEvent is a single composition change
type CompositionChangedEvent struct {
	Type        string       `json:"type"`
	Composition *Composition `json:"composition"`
}

func (e *CompositionChangedEvent) FromBytes(b []byte) errors.Error {
	if err := json.Unmarshal(b, e); err != nil {
		return errors.NewInternal().SetPath("composition/events.FromBytes").SetCode("UNMARSHAL").SetMessage(err.Error())
	}

	return nil
}

func (e *CompositionChangedEvent) ToBytes() ([]byte, errors.Error) {
	b, err := json.Marshal(e)
	if err != nil {
		return nil, errors.NewInternal().SetPath("composition/events.ToBytes").SetCode("MARSHAL").SetMessage(err.Error())
	}
	return b, nil
}

// CompositionUpdatedAutomatically is an event containing all updated compositions after a dependency update
type CompositionUpdatedAutomaticallyEvent struct {
	Type         string         `json:"type"`
	Compositions []*Composition `json:"compositions"`
}

func (e *CompositionUpdatedAutomaticallyEvent) FromBytes(b []byte) errors.Error {
	if err := json.Unmarshal(b, e); err != nil {
		return errors.NewInternal().SetPath("composition/events.FromBytes").SetCode("UNMARSHAL").SetMessage(err.Error())
	}

	return nil
}

func (e *CompositionUpdatedAutomaticallyEvent) ToBytes() ([]byte, errors.Error) {
	b, err := json.Marshal(e)
	if err != nil {
		return nil, errors.NewInternal().SetPath("composition/events.ToBytes").SetCode("MARSHAL").SetMessage(err.Error())
	}
	return b, nil
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

func (e *ArticleExistsEventRequest) FromBytes(b []byte) errors.Error {
	if err := json.Unmarshal(b, e); err != nil {
		return errors.NewInternal().SetPath("composition/events.FromBytes").SetCode("UNMARSHAL").SetMessage(err.Error())
	}

	return nil
}

func (e *ArticleExistsEventRequest) ToBytes() ([]byte, errors.Error) {
	b, err := json.Marshal(e)
	if err != nil {
		return nil, errors.NewInternal().SetPath("composition/events.ToBytes").SetCode("MARSHAL").SetMessage(err.Error())
	}
	return b, nil
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

func (e *ArticleExistsEventResponse) FromBytes(b []byte) errors.Error {
	if err := json.Unmarshal(b, e); err != nil {
		return errors.NewInternal().SetPath("composition/events.FromBytes").SetCode("UNMARSHAL").SetMessage(err.Error())
	}

	return nil
}

func (e *ArticleExistsEventResponse) ToBytes() ([]byte, errors.Error) {
	b, err := json.Marshal(e)
	if err != nil {
		return nil, errors.NewInternal().SetPath("composition/events.ToBytes").SetCode("MARSHAL").SetMessage(err.Error())
	}
	return b, nil
}
