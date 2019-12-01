package user

import (
	"testing"

	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/tests"
)

func hasErrCode(err errors.Error, code string) bool {
	if err == nil {
		return false
	}
	return err.Code() == code
}

func TestValidateSchema(t *testing.T) {
	user := NewUser()

	tests.ErrCode(t, user.ValidateSchema(), "INVALID_USERNAME_LENGTH")

	user.Username = "username"
	tests.ErrCode(t, user.ValidateSchema(), "INVALID_NAME_LENGTH")

	user.Name = "Name"
	tests.ErrCode(t, user.ValidateSchema(), "INVALID_LASTNAME_LENGTH")

	user.Lastname = "Name"
	tests.ErrCode(t, user.ValidateSchema(), "INVALID_EMAIL_LENGTH")

	user.Email = "asd&asd.com"
	tests.ErrCode(t, user.ValidateSchema(), "INVALID_EMAIL_ADDRESS")
	user.Email = "as-as-as"
	tests.ErrCode(t, user.ValidateSchema(), "INVALID_EMAIL_ADDRESS")
	user.Email = "asd@google_yahoo.com"
	tests.ErrCode(t, user.ValidateSchema(), "INVALID_EMAIL_ADDRESS")

	user.Email = "asd@asd.com"

	tests.Ok(t, user.ValidateSchema())
}

func TestUserPassword(t *testing.T) {
	user := NewUser()
	pwd := "123456"
	user.SetPassword(pwd)

	tests.Assert(t, len(user.Password) > 30)
	tests.Assert(t, user.Password != pwd)

	tests.Assert(t, user.ComparePassword("123456"))
	tests.Assert(t, !user.ComparePassword("123457"))
}

func TestUserRoles(t *testing.T) {
	user := NewUser()

	tests.Assert(t, user.HasRole("user"))
	tests.Equal(t, len(user.Roles), 1)

	user.AddRole("admin")
	tests.Assert(t, user.HasRole("admin"))
	tests.Equal(t, len(user.Roles), 2)

	user.AddRole("user")
	tests.Assert(t, user.HasRole("user"))
	tests.Equal(t, len(user.Roles), 2)

	user.RemoveRole("admin")
	tests.Assert(t, !user.HasRole("admin"))
	tests.Assert(t, user.HasRole("user"))
	tests.Equal(t, len(user.Roles), 1)
}
