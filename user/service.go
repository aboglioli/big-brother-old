package user

import (
	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/events"
)

type Service interface {
	GetByID(id string) (*User, errors.Error)
	Create(req *CreateRequest) (*User, errors.Error)
	Update(id string, req *UpdateRequest) (*User, errors.Error)
	Delete(id string) errors.Error
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

func (s *service) GetByID(id string) (*User, errors.Error) {
	user, err := s.repository.FindByID(id)
	if err != nil {
		return nil, errors.NewStatus("USER_NOT_FOUND").SetPath("user/service.GetByID").SetRef(err)
	}
	if !user.Enabled {
		return nil, errors.NewStatus("USER_IS_DELETED").SetPath("user/service.GetByID")
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

func (s *service) Create(req *CreateRequest) (*User, errors.Error) {
	path := "user/service.Create"
	u := NewUser()

	if existing, err := s.repository.FindByUsername(req.Username); existing != nil || err == nil {
		return nil, errors.NewStatus("EXISTING_USERNAME").SetPath(path)
	}

	if existing, err := s.repository.FindByEmail(req.Email); existing != nil || err == nil {
		return nil, errors.NewStatus("EXISTING_EMAIL").SetPath(path)
	}

	if len(req.Password) < 8 {
		return nil, errors.NewStatus("PASSWORD_TOO_SHORT").SetPath(path)
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

	// Publish event: composition.created
	// event, opts := NewCompositionCreatedEvent(c)
	// if err := s.eventMgr.Publish(event, opts); err != nil {
	// 	return nil, err
	// }

	return u, nil
}

type UpdateRequest struct {
	Username *string `json:"username" bson:"username"`
	Password *string `json:"password" bson:"password"`
	Name     *string `json:"name" bson:"name"`
	Lastname *string `json:"lastname" bson:"lastname"`
	Email    *string `json:"email" bson:"email"`
}

func (s *service) Update(id string, req *UpdateRequest) (*User, errors.Error) {
	path := "user/service.Update"

	u, err := s.repository.FindByID(id)
	if u == nil || err != nil {
		return nil, errors.NewStatus("USER_DOES_NOT_EXIST").SetPath(path).SetRef(err)
	}

	if req.Username != nil {
		if existing, err := s.repository.FindByUsername(*req.Username); existing != nil || err == nil {
			return nil, errors.NewStatus("EXISTING_USERNAME").SetPath(path).SetMessage(err.Error())
		}
		u.Username = *req.Username
	}

	if req.Password != nil {
		if len(*req.Password) < 8 {
			return nil, errors.NewStatus("PASSWORD_TOO_SHORT").SetPath(path)
		}
		u.SetPassword(*req.Password)
	}

	if req.Name != nil {
		u.Name = *req.Name
	}

	if req.Lastname != nil {
		u.Lastname = *req.Lastname
	}

	if req.Email != nil {
		if existing, err := s.repository.FindByEmail(*req.Email); existing != nil || err == nil {
			return nil, errors.NewStatus("EXISTING_EMAIL").SetPath(path).SetMessage(err.Error())
		}
		u.Email = *req.Email
	}

	if err := u.ValidateSchema(); err != nil {
		return nil, err
	}

	if err := s.repository.Update(u); err != nil {
		return nil, errors.NewInternal("UPDATE").SetPath(path).SetRef(err)
	}

	return u, nil
}

func (s *service) Delete(id string) errors.Error {
	path := "user/service.Delete"

	_, err := s.repository.FindByID(id)
	if err != nil {
		return errors.NewStatus("NOT_FOUND").SetPath(path).SetRef(err)
	}

	if err := s.repository.Delete(id); err != nil {
		return errors.NewInternal("DELETE").SetPath(path).SetRef(err)
	}

	return nil
}
