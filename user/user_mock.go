package user

func newUser() *User {
	u := NewUser()
	u.Username = "test-user"
	u.SetPassword("12345678")
	u.Name = "Name"
	u.Lastname = "Lastname"
	u.Email = "test@user.com"
	return u
}
