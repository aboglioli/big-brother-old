package user

import (
	"testing"

	"github.com/aboglioli/big-brother/pkg/errors"
)

func hasErrCode(err errors.Error, code string) bool {
	if err == nil {
		return false
	}
	return err.Code() == code
}

func TestUserValidation(t *testing.T) {
	user := NewUser()

	if err := user.ValidateSchema(); !hasErrCode(err, "INVALID_USERNAME_LENGTH") {
		t.Error(err)
	}

	user.Username = "username"

	if err := user.ValidateSchema(); !hasErrCode(err, "INVALID_NAME_LENGTH") {
		t.Error(err)
	}

	user.Name = "Name"

	if err := user.ValidateSchema(); !hasErrCode(err, "INVALID_LASTNAME_LENGTH") {
		t.Error(err)
	}

	user.Lastname = "Name"

	if err := user.ValidateSchema(); !hasErrCode(err, "INVALID_EMAIL_LENGTH") {
		t.Error(err)
	}

	user.Email = "asd&asd.com"
	if err := user.ValidateSchema(); !hasErrCode(err, "INVALID_EMAIL_ADDRESS") {
		t.Error(err)
	}
	user.Email = "as-as-as"
	if err := user.ValidateSchema(); !hasErrCode(err, "INVALID_EMAIL_ADDRESS") {
		t.Error(err)
	}
	user.Email = "asd@google_yahoo.com"
	if err := user.ValidateSchema(); !hasErrCode(err, "INVALID_EMAIL_ADDRESS") {
		t.Error(err)
	}

	user.Email = "asd@asd.com"

	if err := user.ValidateSchema(); err != nil {
		t.Error(err)
	}
}
