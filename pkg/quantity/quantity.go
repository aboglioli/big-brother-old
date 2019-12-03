package quantity

import (
	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/unit"
)

// Quantity defines quantity with unit from International System of Units
type Quantity struct {
	Quantity float64 `bson:"quantity" json:"quantity"`
	Unit     string  `bson:"unit" json:"unit"`
}

func (q1 Quantity) Add(q2 Quantity) (Quantity, error) {
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

func (q1 Quantity) Subtract(q2 Quantity) (Quantity, error) {
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

	if u1 == nil || u2 == nil {
		return false
	}

	n1 := q1.Normalize()
	n2 := q2.Normalize()

	return u1.SameType(u2) && n1 == n2
}

func (q1 Quantity) Compatible(q2 Quantity) bool {
	repo := unit.GetRepository()
	u1, u2 := repo.FindByName(q1.Unit), repo.FindByName(q2.Unit)

	if u1 == nil || u2 == nil {
		return false
	}

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

func referenceUnit(q1, q2 Quantity) (*unit.Unit, error) {
	path := "quantity/quantity.referenceUnit"
	repo := unit.GetRepository()

	u1 := repo.FindByName(q1.Unit)
	if u1 == nil {
		return nil, errors.NewStatus("UNIT_DOES_NOT_EXIST").SetPath(path)
	}

	u2 := repo.FindByName(q2.Unit)
	if u2 == nil {
		return nil, errors.NewStatus("UNIT_DOES_NOT_EXIST").SetPath(path)
	}

	if !u1.SameType(u2) {
		return nil, errors.NewStatus("INCOMPATIBLE_UNITS").SetPath(path)
	}

	return u1, nil
}
