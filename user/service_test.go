package user

import (
	"testing"

	"github.com/aboglioli/big-brother/pkg/events"
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
	repo, eventMgr := newMockRepository(), events.GetMockManager()
	serv := NewService(repo, eventMgr)

	t.Run("Default values", func(t *testing.T) {
		repo.Clean()
		eventMgr.Clean()
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
		assert.Equal(t, createdUser.ID.Hex(), insertedUser.ID.Hex())

		assert.Assert(t, createdUser.ComparePassword(req.Password))
	})

	t.Run("Invalid fields", func(t *testing.T) {
		repo.Clean()
		eventMgr.Clean()
		user, user2 := newUser(), newUser()
		repo.Insert(user2)
		repo.Reset()

		// Username
		req := userToCreateRequest(user)

		createdUser, err := serv.Create(req)
		assert.ErrCode(t, err, "EXISTING_USERNAME")
		assert.Nil(t, createdUser)

		repo.Assert(t, []mock.Call{
			mock.Call{"FindByUsername", []interface{}{user.Username}},
		})

		// Email
		repo.Reset()
		req = userToCreateRequest(user)
		req.Username = "another-user"
		req.Password = "12345678"

		createdUser, err = serv.Create(req)
		assert.ErrCode(t, err, "EXISTING_EMAIL")
		assert.Nil(t, createdUser)

		repo.Assert(t, []mock.Call{
			mock.Call{"FindByUsername", []interface{}{req.Username}},
			mock.Call{"FindByEmail", []interface{}{req.Email}},
		})

		// Password
		repo.Reset()
		req = userToCreateRequest(user)
		req.Username = "another-user"
		req.Email = "another@email.com"
		req.Password = "123456"

		createdUser, err = serv.Create(req)
		assert.ErrCode(t, err, "PASSWORD_TOO_SHORT")
		assert.Nil(t, createdUser)

		repo.Assert(t, []mock.Call{
			mock.Call{"FindByUsername", []interface{}{req.Username}},
			mock.Call{"FindByEmail", []interface{}{req.Email}},
		})

		// ValidateSchema
		repo.Clean()
		user = newUser()
		req = userToCreateRequest(user)
		req.Username = ""
		createdUser, err = serv.Create(req)
		assert.ErrValidation(t, err, "username", "INVALID_LENGTH")
		assert.Nil(t, createdUser)

		req = userToCreateRequest(user)
		req.Name = ""
		createdUser, err = serv.Create(req)
		assert.ErrValidation(t, err, "name", "INVALID_LENGTH")
		assert.Nil(t, createdUser)

		req = userToCreateRequest(user)
		req.Lastname = ""
		createdUser, err = serv.Create(req)
		assert.ErrValidation(t, err, "lastname", "INVALID_LENGTH")
		assert.Nil(t, createdUser)
	})
}
