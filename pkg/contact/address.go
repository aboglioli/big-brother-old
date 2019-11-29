package contact

type Address struct {
	Address   string  `json:"address" bson:"address"`
	Country   string  `json:"country" bson:"country"`
	State     string  `json:"state" bson:"state"`
	ZIPCode   string  `json:"zipCode" bson:"zipCode"`
	Latitude  float64 `json:"lat" bson:"lat"`
	Longitude float64 `json:"lng" bson:"lng"`
}

// TODO: implement
func (a Address) IsValid() bool {
	return true
}
