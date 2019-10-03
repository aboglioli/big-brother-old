package quantity

type Quantity struct {
	Unit     string  `bson:"unit" json:"unit" binding:"required"`
	Quantity float64 `bson:"quantity" json:"value" binding:"required"`
}
