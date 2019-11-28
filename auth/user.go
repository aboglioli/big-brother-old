package auth

type User struct {
	ID          string   `json:"id" bson:"_id"`
	Username    string   `json:"username"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
	Login       string   `json:"login"`
}

func (u *User) HasPermission(perm string) bool {
	for _, p := range u.Permissions {
		if p == perm {
			return true
		}
	}
	return false
}
