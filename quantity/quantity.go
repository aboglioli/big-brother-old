package quantity

type Quantity struct {
	Unit  string  `bson:"unit" json:"unit" binding:"required"`
	Value float32 `bson:"value" json:"value" binding:"required"`
}
