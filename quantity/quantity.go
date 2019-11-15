package quantity

import (
	"github.com/aboglioli/big-brother/errors"
	"github.com/aboglioli/big-brother/unit"
)

type Quantity struct {
	Quantity float64 `bson:"quantity" json:"quantity" validate:"required"`
	Unit     string  `bson:"unit" json:"unit" validate:"required"`
}

func IsValid(q Quantity) bool {
	repo := unit.GetRepository()
	if repo.Exists(q.Unit) && q.Quantity >= 0 {
		return true
	}
	return false
}

func (q1 Quantity) Add(q2 Quantity) (Quantity, errors.Error) {
	u1, err := referenceUnit(q1, q2)
	if err != nil {
		return Quantity{}, err
	}

	nQ1 := q1.Normalize()
	nQ2 := q2.Normalize()
	total := nQ1 + nQ2
	total = total / u1.Modifier

	return Quantity{
		Unit:     q1.Unit,
		Quantity: total,
	}, nil
}

func (q1 Quantity) Subtract(q2 Quantity) (Quantity, errors.Error) {
	u1, err := referenceUnit(q1, q2)
	if err != nil {
		return Quantity{}, err
	}

	nQ1 := q1.Normalize()
	nQ2 := q2.Normalize()
	total := nQ1 - nQ2
	total = total / u1.Modifier

	return Quantity{
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

func (q Quantity) IsValid() bool {
	repo := unit.GetRepository()
	if repo.Exists(q.Unit) && q.Quantity >= 0 {
		return true
	}
	return false
}

func (q Quantity) IsEmpty() bool {
	return q.Quantity == 0 && q.Unit == ""
}

func referenceUnit(q1, q2 Quantity) (*unit.Unit, errors.Error) {
	errGen := errors.NewValidation().SetPath("quantity/quantity.referenceUnit")
	repo := unit.GetRepository()

	u1 := repo.FindByName(q1.Unit)
	if u1 == nil {
		return nil, errGen.SetCode("UNIT_DOES_NOT_EXIST")
	}

	u2 := repo.FindByName(q2.Unit)
	if u2 == nil {
		return nil, errGen.SetCode("UNIT_DOES_NOT_EXIST")
	}

	if !u1.SameType(u2) {
		return nil, errGen.SetCode("INCOMPATIBLE_UNITS")
	}

	return u1, nil
}
