package user

import (
	"testing"

	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/tests/assert"
)

func hasErrCode(err errors.Error, code string) bool {
	if err == nil {
		return false
	}
	return err.Code() == code
}

func TestValidateSchema(t *testing.T) {
	user := NewUser()

	assert.ErrCode(t, user.ValidateSchema(), "INVALID_USERNAME_LENGTH")

	user.Username = "username"
	assert.ErrCode(t, user.ValidateSchema(), "INVALID_NAME_LENGTH")

	user.Name = "Name"
	assert.ErrCode(t, user.ValidateSchema(), "INVALID_LASTNAME_LENGTH")

	user.Lastname = "Name"
	assert.ErrCode(t, user.ValidateSchema(), "INVALID_EMAIL_LENGTH")

	user.Email = "asd&asd.com"
	assert.ErrCode(t, user.ValidateSchema(), "INVALID_EMAIL_ADDRESS")
	user.Email = "as-as-as"
	assert.ErrCode(t, user.ValidateSchema(), "INVALID_EMAIL_ADDRESS")
	user.Email = "asd@google_yahoo.com"
	assert.ErrCode(t, user.ValidateSchema(), "INVALID_EMAIL_ADDRESS")

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
