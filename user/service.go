package user

import (
	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/events"
)

type Service interface {
	GetByID(id string) (*User, error)
	Create(req *CreateRequest) (*User, error)
	Update(id string, req *UpdateRequest) (*User, error)
	Delete(id string) error
}

type service struct {
	repository Repository
	eventMgr   events.Manager
}

func NewService(repo Repository, eventMgr events.Manager) Service {
	return &service{
		repository: repo,
		eventMgr:   eventMgr,
	}
}

func (s *service) GetByID(id string) (*User, error) {
	user, err := s.repository.FindByID(id)
	if err != nil || !user.Enabled {
		return nil, errors.NewStatus("USER_NOT_FOUND").SetPath("user/service.GetByID").SetStatus(404).SetRef(err)
	}
	return user, nil
}

type CreateRequest struct {
	Username string `json:"username" bson:"username" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
	Name     string `json:"name" bson:"name" binding:"required"`
	Lastname string `json:"lastname" bson:"lastname" binding:"required"`
	Email    string `json:"email" bson:"email" binding:"required"`
}

func (s *service) Create(req *CreateRequest) (*User, error) {
	path := "user/service.Create"
	u := NewUser()

	validErr := errors.NewValidation("VALIDATION").SetPath(path)

	if existing, err := s.repository.FindByUsername(req.Username); existing != nil || err == nil {
		validErr.Add("username", "NOT_AVAILABLE")
	}

	if existing, err := s.repository.FindByEmail(req.Email); existing != nil || err == nil {
		validErr.Add("email", "NOT_AVAILABLE")
	}

	if len(req.Password) < 8 {
		validErr.Add("password", "TOO_SHORT")
	}

	if validErr.Size() > 0 {
		return nil, validErr
	}

	u.Username = req.Username
	u.Name = req.Name
	u.Lastname = req.Lastname
	u.Email = req.Email
	u.SetPassword(req.Password)

	if err := u.ValidateSchema(); err != nil {
		return nil, err
	}

	if err := s.repository.Insert(u); err != nil {
		return nil, errors.NewInternal("INSERT").SetPath(path).SetRef(err)
	}

	// Publish event: user.created
	event, opts := NewUserCreatedEvent(u)
	if err := s.eventMgr.Publish(event, opts); err != nil {
		return nil, err
	}

	return u, nil
}

type UpdateRequest struct {
	Username *string `json:"username" bson:"username"`
	Password *string `json:"password" bson:"password"`
	Name     *string `json:"name" bson:"name"`
	Lastname *string `json:"lastname" bson:"lastname"`
	Email    *string `json:"email" bson:"email"`
}

func (s *service) Update(id string, req *UpdateRequest) (*User, error) {
	path := "user/service.Update"

	u, err := s.repository.FindByID(id)
	if u == nil || err != nil {
		return nil, errors.NewStatus("USER_NOT_FOUND").SetPath(path).SetRef(err)
	}

	validErr := errors.NewValidation("VALIDATION").SetPath(path)

	if req.Username != nil {
		if existing, _ := s.repository.FindByUsername(*req.Username); existing != nil && existing.ID.Hex() != u.ID.Hex() {
			validErr.Add("username", "NOT_AVAILABLE")
		}
		u.Username = *req.Username
	}

	if req.Email != nil {
		if existing, _ := s.repository.FindByEmail(*req.Email); existing != nil && existing.ID.Hex() != u.ID.Hex() {
			validErr.Add("email", "NOT_AVAILABLE")
		}
		u.Email = *req.Email
	}

	if req.Password != nil {
		if len(*req.Password) < 8 {
			validErr.Add("password", "TOO_SHORT")
		}
		u.SetPassword(*req.Password)
	}

	if validErr.Size() > 0 {
		return nil, validErr
	}

	if req.Name != nil {
		u.Name = *req.Name
	}

	if req.Lastname != nil {
		u.Lastname = *req.Lastname
	}

	if err := u.ValidateSchema(); err != nil {
		return nil, err
	}

	if err := s.repository.Update(u); err != nil {
		return nil, errors.NewInternal("UPDATE").SetPath(path).SetRef(err)
	}

	// Publish event: user.updated
	event, opts := NewUserUpdatedEvent(u)
	if err := s.eventMgr.Publish(event, opts); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *service) Delete(id string) error {
	path := "user/service.Delete"

	u, err := s.repository.FindByID(id)
	if err != nil {
		return errors.NewStatus("NOT_FOUND").SetPath(path).SetRef(err)
	}

	if err := s.repository.Delete(id); err != nil {
		return errors.NewInternal("DELETE").SetPath(path).SetRef(err)
	}

	u.Enabled = false
	// Event
	event, opts := NewUserDeletedEvent(u)
	if err := s.eventMgr.Publish(event, opts); err != nil {
		return err
	}

	return nil
}
