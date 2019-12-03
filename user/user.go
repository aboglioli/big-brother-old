package user

import (
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

func (u *User) SetPassword(pwd string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return errors.NewStatus("SET_PASSWORD").SetPath("user/user.SetPassword").SetRef(err)
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

func (u *User) ValidateSchema() error {
	err := errors.NewValidation("VALIDATE_SCHEMA").SetPath("user/user.ValidateSchema")

	if len(u.Username) < 6 || len(u.Username) > 64 {
		err.AddWithMessage("username", "INVALID_LENGTH", "%d", len(u.Username))
	}

	if len(u.Name) < 1 || len(u.Name) > 64 {
		err.AddWithMessage("name", "INVALID_LENGTH", "%d", len(u.Name))
	}

	if len(u.Lastname) < 1 || len(u.Lastname) > 64 {
		err.AddWithMessage("lastname", "INVALID_LENGTH", "%d", len(u.Lastname))
	}

	if len(u.Email) < 6 || len(u.Email) > 80 {
		err.AddWithMessage("email", "INVALID_LENGTH", "%d", len(u.Email))
	}

	if !re.MatchString(u.Email) {
		err.AddWithMessage("email", "INVALID_ADDRESS", "%s", u.Email)
	}

	if !u.Address.IsValid() {
		err.Add("address", "INVALID")
	}

	if !u.Contact.IsValid() {
		err.Add("contact", "INVALID")
	}

	if !u.Social.IsValid() {
		err.Add("social", "INVALID")
	}

	if err.Size() > 0 {
		return err
	}

	return nil
}
