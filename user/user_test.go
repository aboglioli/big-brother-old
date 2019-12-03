package user

import (
	"testing"

	"github.com/aboglioli/big-brother/pkg/tests/assert"
)

func TestValidateSchema(t *testing.T) {
	user := NewUser()

	assert.ErrValidation(t, user.ValidateSchema(), "username", "INVALID_LENGTH")
	assert.ErrValidation(t, user.ValidateSchema(), "name", "INVALID_LENGTH")
	assert.ErrValidation(t, user.ValidateSchema(), "lastname", "INVALID_LENGTH")
	assert.ErrValidation(t, user.ValidateSchema(), "email", "INVALID_LENGTH")
	assert.ErrValidation(t, user.ValidateSchema(), "email", "INVALID_ADDRESS")

	user.Username = "username"
	assert.ErrValidation(t, user.ValidateSchema(), "name", "INVALID_LENGTH")

	user.Name = "Name"
	assert.ErrValidation(t, user.ValidateSchema(), "lastname", "INVALID_LENGTH")

	user.Lastname = "Name"
	assert.ErrValidation(t, user.ValidateSchema(), "email", "INVALID_LENGTH")

	user.Email = "asd&asd.com"
	assert.ErrValidation(t, user.ValidateSchema(), "email", "INVALID_ADDRESS")
	user.Email = "as-as-as"
	assert.ErrValidation(t, user.ValidateSchema(), "email", "INVALID_ADDRESS")
	user.Email = "asd@google_yahoo.com"
	assert.ErrValidation(t, user.ValidateSchema(), "email", "INVALID_ADDRESS")

	user.Email = "asd@asd.com"

	assert.Ok(t, user.ValidateSchema())
}

func TestUserPassword(t *testing.T) {
	user := NewUser()
	pwd := "123456"
	user.SetPassword(pwd)

	assert.Assert(t, len(user.Password) > 30)
	assert.Assert(t, user.Password != pwd)

	assert.Assert(t, user.ComparePassword("123456"))
	assert.Assert(t, !user.ComparePassword("123457"))
}

func TestUserRoles(t *testing.T) {
	user := NewUser()

	assert.Assert(t, user.HasRole("user"))
	assert.Equal(t, len(user.Roles), 1)

	user.AddRole("admin")
	assert.Assert(t, user.HasRole("admin"))
	assert.Equal(t, len(user.Roles), 2)

	user.AddRole("user")
	assert.Assert(t, user.HasRole("user"))
	assert.Equal(t, len(user.Roles), 2)

	user.RemoveRole("admin")
	assert.Assert(t, !user.HasRole("admin"))
	assert.Assert(t, user.HasRole("user"))
	assert.Equal(t, len(user.Roles), 1)
}
