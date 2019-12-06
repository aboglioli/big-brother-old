package composition

import (
	"time"

	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/quantity"
	"github.com/aboglioli/big-brother/pkg/tests/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Composition
func newComposition() *Composition {
	comp := NewComposition()
	comp.Unit.Unit = "u"
	comp.Stock.Unit = "u"
	comp.Validated = true
	return comp
}

func makeMockedCompositions() []*Composition {
	p1 := &Composition{
		ID:   primitive.NewObjectID(),
		Cost: 200.0,
		Unit: quantity.Quantity{
			Quantity: 2.0,
			Unit:     "kg",
		},
	}
	p2 := &Composition{
		ID:   primitive.NewObjectID(),
		Cost: 0.0, // 0.2 * 200 / 2 = 20
		Unit: quantity.Quantity{
			Quantity: 0.2,
			Unit:     "kg",
		},
		Dependencies: []Dependency{
			Dependency{
				On: p1.ID,
				Quantity: quantity.Quantity{
					Quantity: 200.0,
					Unit:     "g",
				},
			},
		},
	}
	p3 := &Composition{
		ID:   primitive.NewObjectID(),
		Cost: 0.0, // 0.1 * 200 / 2 = 10
		Unit: quantity.Quantity{
			Quantity: 500.0,
			Unit:     "g",
		},
		Dependencies: []Dependency{
			Dependency{
				On: p1.ID,
				Quantity: quantity.Quantity{
					Quantity: 100.0,
					Unit:     "g",
				},
			},
		},
	}
	p4 := &Composition{
		ID:   primitive.NewObjectID(),
		Cost: 150.0,
		Unit: quantity.Quantity{
			Quantity: 100.0,
			Unit:     "g",
		},
	}
	p5 := &Composition{
		ID:   primitive.NewObjectID(),
		Cost: 0.0, // 0.4*20/0.2 + 0.05*10/0.5 = 41
		Unit: quantity.Quantity{
			Quantity: 1.0,
			Unit:     "u",
		},
		Dependencies: []Dependency{
			Dependency{
				On: p2.ID,
				Quantity: quantity.Quantity{
					Quantity: 400.0,
					Unit:     "g",
				},
			},
			Dependency{
				On: p3.ID,
				Quantity: quantity.Quantity{
					Quantity: 50.0,
					Unit:     "g",
				},
			},
		},
	}
	p6 := &Composition{
		ID:   primitive.NewObjectID(),
		Cost: 0.0, // 0.35*150/0.1 = 525
		Unit: quantity.Quantity{
			Quantity: 2.0,
			Unit:     "u",
		},
		Dependencies: []Dependency{
			Dependency{
				On: p4.ID,
				Quantity: quantity.Quantity{
					Quantity: 350.0,
					Unit:     "g",
				},
			},
		},
	}
	p7 := &Composition{
		ID:   primitive.NewObjectID(),
		Cost: 0.0, // 2*41/1 + 1.5*525/2 = 475.75
		Unit: quantity.Quantity{
			Quantity: 3.0,
			Unit:     "u",
		},
		Dependencies: []Dependency{
			Dependency{
				On: p5.ID,
				Quantity: quantity.Quantity{
					Quantity: 2.0,
					Unit:     "u",
				},
			},
			Dependency{
				On: p6.ID,
				Quantity: quantity.Quantity{
					Quantity: 1.5,
					Unit:     "u",
				},
			},
		},
	}
	comps := []*Composition{p1, p2, p3, p4, p5, p6, p7}

	for _, c := range comps {
		c.AutoupdateCost = true
		c.Enabled = true
		c.Validated = true
		c.Stock = quantity.Quantity{
			Quantity: 10 * c.Unit.Quantity,
			Unit:     c.Unit.Unit,
		}
	}

	return comps
}

// Repository
type mockRepository struct {
	mock.Mock
	compositions []*Composition
}

func newMockRepository() *mockRepository {
	return &mockRepository{}
}

func (r *mockRepository) Clean() {
	r.compositions = make([]*Composition, 0)
}

func (r *mockRepository) FindAll() ([]*Composition, error) {
	r.Called("FindAll")

	comps := make([]*Composition, 0)
	for _, c := range r.compositions {
		if c.Enabled {
			comps = append(comps, copyComposition(c))
		}
	}

	return comps, nil
}

func (r *mockRepository) FindByID(id string) (*Composition, error) {
	r.Called("FindByID", id)

	for _, c := range r.compositions {
		if c.ID.Hex() == id {
			return copyComposition(c), nil
		}
	}

	return nil, errors.NewInternal("NOT_FOUND").SetPath("composition/repository_mock.FindById")
}

func (r *mockRepository) FindUses(id string) ([]*Composition, error) {
	r.Called("FindUses", id)

	comps := make([]*Composition, 0)
	for _, c := range r.compositions {
		for _, d := range c.Dependencies {
			if d.On.Hex() == id {
				comps = append(comps, copyComposition(c))
				break
			}
		}
	}

	return comps, nil
}

func (r *mockRepository) FindByUsesUpdatedSinceLastChange(usesUpdated bool) ([]*Composition, error) {
	r.Called("FindByUsesUpdatedSinceLastChange", usesUpdated)

	comps := make([]*Composition, 0)
	for _, c := range r.compositions {
		if c.UsesUpdatedSinceLastChange == usesUpdated {
			comps = append(comps, copyComposition(c))
		}
	}

	return comps, nil
}

func (r *mockRepository) Insert(c *Composition) error {
	r.Called("Insert", c)

	c.UpdatedAt = time.Now()
	r.compositions = append(r.compositions, copyComposition(c))

	return nil
}

func (r *mockRepository) InsertMany(comps []*Composition) error {
	r.Called("InsertMany", comps)

	newComps := make([]*Composition, len(comps))
	for i, c := range comps {
		c.UpdatedAt = time.Now()
		newComps[i] = copyComposition(c)
	}
	r.compositions = append(r.compositions, newComps...)

	return nil
}

func (r *mockRepository) Update(c *Composition) error {
	r.Called("Update", c)

	for _, comp := range r.compositions {
		if comp.ID.Hex() == c.ID.Hex() {
			*comp = *copyComposition(c)
			comp.UpdatedAt = time.Now()
			break
		}
	}

	return nil
}

func (r *mockRepository) Delete(id string) error {
	r.Called("Delete", id)

	for _, comp := range r.compositions {
		if comp.ID.Hex() == id {
			comp.UpdatedAt = time.Now()
			comp.Enabled = false
			return nil
		}
	}

	return errors.NewInternal("NOT_FOUND").SetPath("composition/repository_mock.Delete")
}

func (r *mockRepository) Count() (int, int) {
	totalCount, enabledCount := 0, 0
	for _, c := range r.compositions {
		totalCount++
		if c.Enabled {
			enabledCount++
		}
	}

	return totalCount, enabledCount
}

func copyComposition(c *Composition) *Composition {
	copy := *c
	return &copy
}
