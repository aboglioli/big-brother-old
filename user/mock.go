package user

import (
	"time"

	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/tests/mock"
)

// User
func newUser() *User {
	u := NewUser()
	u.Username = "test-user"
	u.SetPassword("12345678")
	u.Name = "Name"
	u.Lastname = "Lastname"
	u.Email = "test@user.com"
	return u
}

// Repository
type mockRepository struct {
	mock.Mock
	users []*User
}

func newMockRepository() *mockRepository {
	return &mockRepository{}
}

func (r *mockRepository) Clean() {
	r.users = make([]*User, 0)
}

func (r *mockRepository) FindByID(id string) (*User, error) {
	r.Called(mock.Call("FindByID", id))

	for _, u := range r.users {
		if u.ID.Hex() == id {
			return copyUser(u), nil
		}
	}

	return nil, errors.NewInternal("NOT_FOUND").SetPath("user/repository_mock.FindById")
}

func (r *mockRepository) FindByUsername(username string) (*User, error) {
	r.Called(mock.Call("FindByUsername", username))

	for _, u := range r.users {
		if u.Username == username {
			return copyUser(u), nil
		}
	}

	return nil, errors.NewInternal("NOT_FOUND").SetPath("user/repository_mock.FindByUsername")
}

func (r *mockRepository) FindByEmail(email string) (*User, error) {
	r.Called(mock.Call("FindByEmail", email))

	for _, u := range r.users {
		if u.Email == email {
			return copyUser(u), nil
		}
	}

	return nil, errors.NewInternal("NOT_FOUND").SetPath("user/repository_mock.FindByEmail")
}

func (r *mockRepository) Insert(u *User) error {
	r.Called(mock.Call("Insert", u))

	u.UpdatedAt = time.Now()
	r.users = append(r.users, copyUser(u))

	return nil
}

func (r *mockRepository) InsertMany(users []*User) error {
	r.Called(mock.Call("InsertMany", users))

	newUsers := make([]*User, len(users))
	for i, u := range users {
		u.UpdatedAt = time.Now()
		newUsers[i] = copyUser(u)
	}
	r.users = append(r.users, newUsers...)

	return nil
}

func (r *mockRepository) Update(u *User) error {
	r.Called(mock.Call("Update", u))

	for _, user := range r.users {
		if user.ID.Hex() == u.ID.Hex() {
			*user = *copyUser(u)
			user.UpdatedAt = time.Now()
			break
		}
	}

	return nil
}

func (r *mockRepository) Delete(id string) error {
	r.Called(mock.Call("Delete", id))

	for _, user := range r.users {
		if user.ID.Hex() == id {
			user.UpdatedAt = time.Now()
			user.Enabled = false
			return nil
		}
	}

	return errors.NewInternal("NOT_FOUND").SetPath("user/repository_mock.Delete")
}

func (r *mockRepository) Count() int {
	count := 0
	for _, u := range r.users {
		if u.Enabled {
			count++
		}
	}

	return count
}

func copyUser(u *User) *User {
	copy := *u
	return &copy
}
