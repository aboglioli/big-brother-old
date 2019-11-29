package contact

type Social struct {
	Facebook  string `json:"facebook" bson:"facebook"`
	Twitter   string `json:"twitter" bson:"twitter"`
	Instagram string `json:"instagram" bson:"instagram"`
}

// TODO: implement
func (s Social) IsValid() bool {
	return true
}
