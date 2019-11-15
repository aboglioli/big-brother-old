package auth

type User struct {
	ID          string   `json:"id" validate:"required"`
	Name        string   `json:"name" validate:"required"`
	Permissions []string `json:"permissions"`
	Login       string   `json:"login" validated:"required"`
}

func (u *User) HasPermission(perm string) bool {
	for _, p := range u.Permissions {
		if p == perm {
			return true
		}
	}
	return false
}
