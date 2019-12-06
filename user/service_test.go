package user

import (
	"testing"

	"github.com/aboglioli/big-brother/impl/events"
	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/tests/assert"
	"github.com/aboglioli/big-brother/pkg/tests/mock"
)

func userToCreateRequest(u *User) *CreateRequest {
	return &CreateRequest{
		Username: u.Username,
		Password: "12345678",
		Name:     u.Name,
		Lastname: u.Lastname,
		Email:    u.Email,
	}
}

func TestCreateUser(t *testing.T) {
	repo, eventMgr := newMockRepository(), events.InMemory()
	serv := NewService(repo, eventMgr)

	// Errors
	t.Run("Existing user (username or email)", func(t *testing.T) {
		repo.Clean()
		user := newUser()
		repo.Insert(user)
		repo.Mock.Reset()

		// Same username and email, short password
		req := userToCreateRequest(user)
		req.Password = "1234"
		_, err := serv.Create(req)
		assert.ErrCode(t, err, "VALIDATION")
		assert.ErrValidation(t, err, "username", "NOT_AVAILABLE")
		assert.ErrValidation(t, err, "password", "TOO_SHORT")
		assert.ErrValidation(t, err, "email", "NOT_AVAILABLE")
		repo.Mock.Assert(t, []mock.Call{
			mock.Call{"FindByUsername", []interface{}{req.Username}},
			mock.Call{"FindByEmail", []interface{}{req.Email}},
		})

		// Same username, short password
		repo.Mock.Reset()
		req.Password = "1234567"
		req.Email = "antoher@user.com"
		_, err = serv.Create(req)
		assert.ErrCode(t, err, "VALIDATION")
		assert.ErrValidation(t, err, "username", "NOT_AVAILABLE")
		assert.ErrValidation(t, err, "password", "TOO_SHORT")
		repo.Mock.Assert(t, []mock.Call{
			mock.Call{"FindByUsername", []interface{}{req.Username}},
			mock.Call{"FindByEmail", []interface{}{req.Email}},
		})

		// Same username
		repo.Mock.Reset()
		req.Password = "12345678"
		req.Email = "antoher@user.com"
		_, err = serv.Create(req)
		assert.ErrCode(t, err, "VALIDATION")
		assert.ErrValidation(t, err, "username", "NOT_AVAILABLE")
		repo.Mock.Assert(t, []mock.Call{
			mock.Call{"FindByUsername", []interface{}{req.Username}},
			mock.Call{"FindByEmail", []interface{}{req.Email}},
		})
	})

	t.Run("Validate schema", func(t *testing.T) {
		repo.Clean()
		repo.Mock.Reset()
		req := userToCreateRequest(new(User))
		_, err := serv.Create(req)
		assert.ErrCode(t, err, "SCHEMA")
		val, ok := err.(*errors.Validation)
		assert.Assert(t, ok)
		assert.Equal(t, val.Size(), 5)
	})

	// OK
	t.Run("Default values", func(t *testing.T) {
		repo.Clean()
		repo.Mock.Reset()
		eventMgr.Clean()
		eventMgr.Mock.Reset()
		user := newUser()
		req := userToCreateRequest(user)

		createdUser, err := serv.Create(req)
		assert.Ok(t, err)
		assert.NotNil(t, createdUser)
		assert.Equal(t, user.Username, createdUser.Username)

		repo.Assert(t, []mock.Call{
			mock.Call{"FindByUsername", []interface{}{user.Username}},
			mock.Call{"FindByEmail", []interface{}{user.Email}},
			mock.Call{"Insert", []interface{}{mock.NotNil}},
		})
		insertedUser, ok := repo.Calls[2].Args[0].(*User)
		assert.Assert(t, ok)
		assert.Equal(t, createdUser, insertedUser)
		assert.Equal(t, createdUser.ID.Hex(), insertedUser.ID.Hex())
		assert.Assert(t, createdUser.ComparePassword(req.Password))

		// Default values
		assert.Assert(t, len(createdUser.ID.Hex()) > 0)
		assert.Equal(t, createdUser.Enabled, true)
		assert.Equal(t, createdUser.Active, true)
		assert.Equal(t, createdUser.Validated, false)

		// Events
		eventMgr.Mock.Assert(t, []mock.Call{
			mock.Call{"Publish", []interface{}{mock.NotNil, mock.NotNil}},
		})
		createdEvent, ok := eventMgr.Calls[0].Args[0].(*UserChangedEvent)
		assert.Assert(t, ok)
		assert.Equal(t, createdEvent.User.ID.Hex(), createdUser.ID.Hex())
	})
}

func TestUpdateUser(t *testing.T) {
	repo, eventMgr := newMockRepository(), events.InMemory()
	_ = NewService(repo, eventMgr)
}
