package quantity

import (
	"github.com/aboglioli/big-brother/errors"
	"github.com/aboglioli/big-brother/unit"
)

type Op int

const (
	OP_ADD Op = iota
	OP_SUBSTRACT
)

type Quantity struct {
	Unit     string  `bson:"unit" json:"unit" binding:"required"`
	Quantity float64 `bson:"quantity" json:"value" binding:"required"`
}

func IsValid(q Quantity) bool {
	repo := unit.GetRepository()
	if repo.Exists(q.Unit) && q.Quantity >= 0 {
		return true
	}
	return false
}

func (q1 Quantity) Op(q2 Quantity, op Op) (*Quantity, error) {
	errGen := errors.FromPath("quantity/quantity.Add")
	repo := unit.GetRepository()
	u1 := repo.FindByName(q1.Unit)
	if u1 == nil {
		return nil, errGen("UNIT_DOES_NOT_EXIST", "")
	}

	u2 := repo.FindByName(q2.Unit)
	if u2 == nil {
		return nil, errGen("UNIT_DOES_NOT_EXIST", "")
	}

	if !u1.SameType(u2) {
		return nil, errGen("INCOMPATIBLE_UNITS", "")
	}

	nQ1 := q1.Normalize()
	nQ2 := q2.Normalize()

	var total float64
	if op == OP_ADD {
		total = nQ1 + nQ2
	} else {
		total = nQ1 - nQ2
	}
	total = total / u1.Modifier

	return &Quantity{
		Unit:     q1.Unit,
		Quantity: total,
	}, nil
}

func (q1 Quantity) Equals(q2 Quantity) bool {
	repo := unit.GetRepository()

	u1, u2 := repo.FindByName(q1.Unit), repo.FindByName(q2.Unit)

	n1 := q1.Normalize()
	n2 := q2.Normalize()

	return u1.SameType(u2) && n1 == n2
}

func (q1 Quantity) Compatible(q2 Quantity) bool {
	repo := unit.GetRepository()
	u1, u2 := repo.FindByName(q1.Unit), repo.FindByName(q2.Unit)
	return u1.SameType(u2)
}

func (q Quantity) Normalize() float64 {
	repo := unit.GetRepository()
	u := repo.FindByName(q.Unit)
	return q.Quantity * u.Modifier
}
