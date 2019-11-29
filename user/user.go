package user

import (
	"fmt"
	"regexp"
	"time"

	"github.com/aboglioli/big-brother/pkg/contact"
	"github.com/aboglioli/big-brother/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Username string             `json:"username" bson:"username"`
	Password string             `json:"-" bson:"password"`
	Name     string             `json:"name" bson:"name"`
	Lastname string             `json:"lastname" bson:"lastname"`
	Email    string             `json:"email" bson:"email"`
	Roles    []string           `json:"roles" bson:"roles"`

	Address contact.Address `json:"address" bson:"address"`
	Contact contact.Contact `json:"contact" bson:"contact"`
	Social  contact.Social  `json:"social" bson:"social"`

	Enabled   bool      `json:"enabled" bson:"enabled"`
	Validated bool      `json:"validated" bson:"validated"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

func NewUser() *User {
	return &User{
		ID:        primitive.NewObjectID(),
		Enabled:   true,
		Validated: false,
		Roles:     []string{"user"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (u *User) SetPassword(pwd string) errors.Error {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return errors.NewValidation().SetPath("auth/user.SetPassword").SetCode("SET_PASSWORD").SetMessage(err.Error())
	}
	u.Password = string(hash)
	return nil
}

func (u *User) ComparePassword(pwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pwd)); err != nil {
		return false
	}
	return true
}

func (u *User) HasRole(role string) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

func (u *User) AddRole(role string) {
	if !u.HasRole(role) {
		u.Roles = append(u.Roles, role)
	}
}

func (u *User) RemoveRole(role string) {
	if u.HasRole(role) {
		roles := make([]string, 0)
		for _, r := range u.Roles {
			if r != role {
				roles = append(roles, r)
			}
		}
		u.Roles = roles
	}
}

var re = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func (u *User) ValidateSchema() errors.Error {
	errGen := errors.NewValidation().SetPath("auth/user.ValidateSchema")

	if len(u.Username) < 6 || len(u.Username) > 64 {
		return errGen.SetCode("INVALID_USERNAME_LENGTH").SetMessage(fmt.Sprintf("%d", len(u.Username)))
	}

	if len(u.Name) < 1 || len(u.Name) > 64 {
		return errGen.SetCode("INVALID_NAME_LENGTH").SetMessage(fmt.Sprintf("%d", len(u.Name)))
	}

	if len(u.Lastname) < 1 || len(u.Lastname) > 64 {
		return errGen.SetCode("INVALID_LASTNAME_LENGTH").SetMessage(fmt.Sprintf("%d", len(u.Lastname)))
	}

	if len(u.Email) < 6 || len(u.Email) > 64 {
		return errGen.SetCode("INVALID_EMAIL_LENGTH").SetMessage(fmt.Sprintf("%d", len(u.Email)))
	}

	if !re.MatchString(u.Email) {
		return errGen.SetCode("INVALID_EMAIL_ADDRESS")
	}

	if !u.Address.IsValid() {
		return errGen.SetCode("INVALID_ADDRESS")
	}

	if !u.Contact.IsValid() {
		return errGen.SetCode("INVALID_CONTACT")
	}

	if !u.Social.IsValid() {
		return errGen.SetCode("INVALID_SOCIAL")
	}

	return nil
}
