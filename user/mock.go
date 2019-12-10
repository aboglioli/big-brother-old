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
	call := mock.Call("FindByID", id)

	for _, u := range r.users {
		if u.ID.Hex() == id {
			r.Mock.Called(call.Return(copyUser(u), nil))
			return copyUser(u), nil
		}
	}

	err := errors.NewInternal("NOT_FOUND").SetPath("user/mock.FindById")
	r.Mock.Called(call.Return(nil, err))
	return nil, err
}

func (r *mockRepository) FindByUsername(username string) (*User, error) {
	call := mock.Call("FindByUsername", username)

	for _, u := range r.users {
		if u.Username == username {
			r.Mock.Called(call.Return(copyUser(u), nil))
			return copyUser(u), nil
		}
	}

	err := errors.NewInternal("NOT_FOUND").SetPath("user/mock.FindByUsername")
	r.Mock.Called(call.Return(nil, err))
	return nil, err
}

func (r *mockRepository) FindByEmail(email string) (*User, error) {
	call := mock.Call("FindByEmail", email)

	for _, u := range r.users {
		if u.Email == email {
			r.Mock.Called(call.Return(copyUser(u), nil))
			return copyUser(u), nil
		}
	}

	err := errors.NewInternal("NOT_FOUND").SetPath("user/mock.FindByEmail")
	r.Mock.Called(call.Return(nil, err))
	return nil, err
}

func (r *mockRepository) Insert(u *User) error {
	call := mock.Call("Insert", u)

	u.UpdatedAt = time.Now()
	r.users = append(r.users, copyUser(u))

	r.Mock.Called(call.Return(nil))
	return nil
}

func (r *mockRepository) InsertMany(users []*User) error {
	call := mock.Call("InsertMany", users)

	newUsers := make([]*User, len(users))
	for i, u := range users {
		u.UpdatedAt = time.Now()
		newUsers[i] = copyUser(u)
	}
	r.users = append(r.users, newUsers...)

	r.Mock.Called(call.Return(nil))
	return nil
}

func (r *mockRepository) Update(u *User) error {
	call := mock.Call("Update", u)

	for _, user := range r.users {
		if user.ID.Hex() == u.ID.Hex() {
			*user = *copyUser(u)
			user.UpdatedAt = time.Now()
			break
		}
	}

	r.Mock.Called(call.Return(nil))
	return nil
}

func (r *mockRepository) Delete(id string) error {
	call := mock.Call("Delete", id)

	for _, user := range r.users {
		if user.ID.Hex() == id {
			user.UpdatedAt = time.Now()
			user.Enabled = false
			r.Mock.Called(call.Return(nil))
			return nil
		}
	}

	err := errors.NewInternal("NOT_FOUND").SetPath("user/mock.Delete")
	r.Mock.Called(call.Return(err))
	return err
}

func (r *mockRepository) Count() (int, int) {
	totalCount, enabledCount := 0, 0
	for _, c := range r.users {
		totalCount++
		if c.Enabled {
			enabledCount++
		}
	}

	return totalCount, enabledCount
}

func copyUser(u *User) *User {
	copy := *u
	return &copy
}
