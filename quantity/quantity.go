package quantity

type Quantity struct {
	Unit     string  `bson:"unit" json:"unit" binding:"required"`
	Quantity float64 `bson:"quantity" json:"value" binding:"required"`
}

func (q1 Quantity) Equals(q2 Quantity) bool {
	return q1.Unit == q2.Unit && q1.Quantity == q2.Quantity
}
