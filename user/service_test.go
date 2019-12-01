package user

import (
	"testing"

	"github.com/aboglioli/big-brother/pkg/events"
	"github.com/aboglioli/big-brother/pkg/tests/assert"
)

func userToCreateRequest(u *User) *CreateRequest {
	return &CreateRequest{
		Username: u.Username,
		Name:     u.Name,
		Lastname: u.Lastname,
		Email:    u.Email,
	}
}

func TestCreateUser(t *testing.T) {
	repo, eventMgr := newMockRepository(), events.GetMockManager()
	serv := NewService(repo, eventMgr)

	t.Run("Default values", func(t *testing.T) {
		user := newUser()
		req := userToCreateRequest(user)
		req.Password = "12345678"

		createdUser, err := serv.Create(req)
		assert.Ok(t, err)
		assert.NotNil(t, createdUser)
	})
}
