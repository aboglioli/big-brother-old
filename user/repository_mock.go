package user

import (
	"time"

	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/tests/mock"
)

type mockRepository struct {
	mock.Mock
	users []*User
}

func newMockRepository() *mockRepository {
	return &mockRepository{}
}

// Helpers
func (r *mockRepository) Clean() {
	r.users = make([]*User, 0)
}

// Implementation
func (r *mockRepository) FindByID(id string) (*User, errors.Error) {
	r.Called("FindByID", id)

	for _, c := range r.users {
		if c.ID.Hex() == id && c.Enabled {
			return copyUser(c), nil
		}
	}

	return nil, errors.NewInternal("NOT_FOUND").SetPath("user/repository_mock.FindById")
}

func (r *mockRepository) FindByUsername(username string) (*User, errors.Error) {
	r.Called("FindByUsername", username)

	for _, c := range r.users {
		if c.Username == username && c.Enabled {
			return copyUser(c), nil
		}
	}

	return nil, errors.NewInternal("NOT_FOUND").SetPath("user/repository_mock.FindByUsername")
}

func (r *mockRepository) FindByEmail(email string) (*User, errors.Error) {
	r.Called("FindByEmail", email)

	for _, c := range r.users {
		if c.Email == email && c.Enabled {
			return copyUser(c), nil
		}
	}

	return nil, errors.NewInternal("NOT_FOUND").SetPath("user/repository_mock.FindByEmail")
}

func (r *mockRepository) Insert(u *User) errors.Error {
	r.Called("Insert", u)

	u.UpdatedAt = time.Now()
	r.users = append(r.users, copyUser(u))

	return nil
}

func (r *mockRepository) InsertMany(users []*User) errors.Error {
	r.Called("InsertMany", users)

	newUsers := make([]*User, len(users))
	for i, u := range users {
		u.UpdatedAt = time.Now()
		newUsers[i] = copyUser(u)
	}
	r.users = append(r.users, newUsers...)

	return nil
}

func (r *mockRepository) Update(u *User) errors.Error {
	r.Called("Update", u)

	for _, user := range r.users {
		if user.ID.Hex() == u.ID.Hex() {
			if !user.Enabled {
				return errors.NewInternal("DISABLED").SetPath("user/repository_mock.Update")
			}
			*user = *copyUser(u)
			user.UpdatedAt = time.Now()
			break
		}
	}

	return nil
}

func (r *mockRepository) Delete(id string) errors.Error {
	r.Called("Delete", id)

	for _, user := range r.users {
		if user.ID.Hex() == id && user.Enabled {
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
