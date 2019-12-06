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

func userToUpdateRequest(u *User) *UpdateRequest {
	password := "12345678"
	return &UpdateRequest{
		Username: &u.Username,
		Password: &password,
		Name:     &u.Name,
		Lastname: &u.Lastname,
		Email:    &u.Email,
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
		assert.Equal(t, err.(*errors.Validation).Size(), 3)
		repo.Mock.Assert(t,
			mock.Call("FindByUsername", req.Username),
			mock.Call("FindByEmail", req.Email),
		)

		// Same username, short password
		repo.Mock.Reset()
		req.Password = "1234567"
		req.Email = "another@user.com"
		_, err = serv.Create(req)
		assert.ErrCode(t, err, "VALIDATION")
		assert.ErrValidation(t, err, "username", "NOT_AVAILABLE")
		assert.ErrValidation(t, err, "password", "TOO_SHORT")
		assert.Equal(t, err.(*errors.Validation).Size(), 2)
		repo.Mock.Assert(t,
			mock.Call("FindByUsername", req.Username),
			mock.Call("FindByEmail", req.Email),
		)

		// Same username
		repo.Mock.Reset()
		req.Password = "12345678"
		req.Email = "antoher@user.com"
		_, err = serv.Create(req)
		assert.ErrCode(t, err, "VALIDATION")
		assert.ErrValidation(t, err, "username", "NOT_AVAILABLE")
		assert.Equal(t, err.(*errors.Validation).Size(), 1)
		repo.Mock.Assert(t,
			mock.Call("FindByUsername", req.Username),
			mock.Call("FindByEmail", req.Email),
		)
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

		repo.Mock.Assert(t,
			mock.Call("FindByUsername", user.Username),
			mock.Call("FindByEmail", user.Email),
			mock.Call("Insert", mock.NotNil),
		)
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
		eventMgr.Mock.Assert(t,
			mock.Call("Publish", mock.NotNil, mock.NotNil),
		)
		event, ok := eventMgr.Calls[0].Args[0].(*UserChangedEvent)
		assert.Assert(t, ok)
		assert.Equal(t, event.Type, "UserCreated")
		assert.Equal(t, event.User.ID.Hex(), createdUser.ID.Hex())
	})
}

func TestUpdateUser(t *testing.T) {
	repo, eventMgr := newMockRepository(), events.InMemory()
	serv := NewService(repo, eventMgr)

	// Errors
	t.Run("Non-existing user", func(t *testing.T) {
		user := newUser()
		_, err := serv.Update("123", userToUpdateRequest(user))
		assert.ErrCode(t, err, "USER_NOT_FOUND")
	})

	t.Run("Existing username or email, short password", func(t *testing.T) {
		user1, user2 := newUser(), newUser()
		user1.Username = "user-1"
		user1.Email = "user-1@user.com"
		user2.Username = "user-2"
		user2.Email = "user-2@user.com"
		repo.InsertMany([]*User{user1, user2})
		repo.Mock.Reset()

		req := userToUpdateRequest(user2)
		username := "user-1"
		email := "user-1@user.com"
		req.Username = &username
		req.Email = &email
		_, err := serv.Update(user2.ID.Hex(), req)
		assert.ErrCode(t, err, "VALIDATION")
		assert.ErrValidation(t, err, "username", "NOT_AVAILABLE")
		assert.ErrValidation(t, err, "email", "NOT_AVAILABLE")
		assert.Equal(t, err.(*errors.Validation).Size(), 2)
		repo.Mock.Assert(t,
			mock.Call("FindByID", user2.ID.Hex()),
			mock.Call("FindByUsername", *req.Username),
			mock.Call("FindByEmail", *req.Email),
		)

		req = userToUpdateRequest(user2)
		email = "user-1@user.com"
		req.Email = &email
		_, err = serv.Update(user2.ID.Hex(), req)
		assert.ErrCode(t, err, "VALIDATION")
		assert.ErrValidation(t, err, "email", "NOT_AVAILABLE")
		assert.Equal(t, err.(*errors.Validation).Size(), 1)

		user2.Username = "another-user"
		user2.Email = "another@user.com"
		req = userToUpdateRequest(user2)
		pwd := "1234567"
		req.Password = &pwd
		_, err = serv.Update(user2.ID.Hex(), req)
		assert.ErrCode(t, err, "VALIDATION")
		assert.ErrValidation(t, err, "password", "TOO_SHORT")
		assert.Equal(t, err.(*errors.Validation).Size(), 1)

	})

	t.Run("Validate schema", func(t *testing.T) {
		repo.Clean()
		user := new(User)
		repo.Insert(user)
		req := userToUpdateRequest(user)
		_, err := serv.Update(user.ID.Hex(), req)
		assert.ErrCode(t, err, "SCHEMA")
		assert.ErrValidation(t, err, "username", "INVALID_LENGTH")
		assert.ErrValidation(t, err, "name", "INVALID_LENGTH")
		assert.ErrValidation(t, err, "lastname", "INVALID_LENGTH")
		assert.ErrValidation(t, err, "email", "INVALID_LENGTH")
		assert.ErrValidation(t, err, "email", "INVALID_ADDRESS")
		assert.Equal(t, err.(*errors.Validation).Size(), 5)
	})

	// OK
	t.Run("Update with default values", func(t *testing.T) {
		repo.Clean()
		user := newUser()
		repo.Insert(user)
		repo.Mock.Reset()
		eventMgr.Clean()
		eventMgr.Mock.Reset()

		req := userToUpdateRequest(user)
		updatedUser, err := serv.Update(user.ID.Hex(), req)
		assert.Ok(t, err)
		assert.Equal(t, updatedUser.ID.Hex(), user.ID.Hex())

		repo.Mock.Assert(t,
			mock.Call("FindByID", user.ID.Hex()),
			mock.Call("FindByUsername", *req.Username),
			mock.Call("FindByEmail", *req.Email),
			mock.Call("Update", mock.NotNil),
		)
		userInDB := repo.Calls[3].Args[0].(*User)
		assert.Equal(t, updatedUser, userInDB)

		eventMgr.Mock.Assert(t,
			mock.Call("Publish", mock.NotNil, mock.NotNil),
		)
		event, ok := eventMgr.Mock.Calls[0].Args[0].(*UserChangedEvent)
		assert.Assert(t, ok)
		assert.Equal(t, event.Type, "UserUpdated")
		assert.Equal(t, updatedUser, event.User)
	})

	t.Run("Change username, email and password", func(t *testing.T) {
		repo.Clean()
		user := newUser()
		repo.Insert(user)
		repo.Mock.Reset()
		eventMgr.Clean()
		eventMgr.Mock.Reset()

		req := userToUpdateRequest(user)
		username := "another-user"
		email := "another@user.com"
		pwd := "88889999"
		req.Username = &username
		req.Email = &email
		req.Password = &pwd
		updatedUser, err := serv.Update(user.ID.Hex(), req)
		assert.Ok(t, err)
		assert.Equal(t, updatedUser.ID.Hex(), user.ID.Hex())
		assert.Equal(t, updatedUser.Username, username)
		assert.Equal(t, updatedUser.Email, email)
		assert.Assert(t, updatedUser.ComparePassword(pwd))

		repo.Mock.Assert(t,
			mock.Call("FindByID", user.ID.Hex()),
			mock.Call("FindByUsername", *req.Username),
			mock.Call("FindByEmail", *req.Email),
			mock.Call("Update", mock.NotNil),
		)
		userInDB := repo.Calls[3].Args[0].(*User)
		assert.Equal(t, updatedUser, userInDB)
	})
}
