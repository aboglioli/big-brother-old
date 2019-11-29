package contact

type Contact struct {
	Email  string `json:"email" bson:"email"`
	Mobile string `json:"mobile" bson:"mobile"`
	Phone  string `json:"phone" bson:"phone"`
}

// TODO: implement
func (c Contact) IsValid() bool {
	return true
}
