package user

import (
	"github.com/aboglioli/big-brother/pkg/events"
)

type UserChangedEvent struct {
	events.Event
	User *User `json:"user"`
}

func NewUserCreatedEvent(u *User) (*UserChangedEvent, *events.Options) {
	event := UserChangedEvent{events.Event{"UserCreated"}, u}
	opts := events.Options{"user", "topic", "user.created", ""}
	return &event, &opts
}

func NewUserUpdatedEvent(u *User) (*UserChangedEvent, *events.Options) {
	event := UserChangedEvent{events.Event{"UserUpdated"}, u}
	opts := events.Options{"user", "topic", "user.updated", ""}
	return &event, &opts
}

func NewUserDeletedEvent(u *User) (*UserChangedEvent, *events.Options) {
	event := UserChangedEvent{events.Event{"UserDeleted"}, u}
	opts := events.Options{"user", "topic", "user.deleted", ""}
	return &event, &opts
}
