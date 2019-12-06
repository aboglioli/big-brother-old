package composition

import (
	"time"

	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/quantity"
	"github.com/aboglioli/big-brother/pkg/tests/mock"
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
	p1 := NewComposition()
	p1.Cost = 200.0
	p1.Unit = quantity.Quantity{2.0, "kg"}

	p2 := NewComposition()
	p2.Unit = quantity.Quantity{0.2, "kg"}
	p2.Dependencies = []Dependency{
		Dependency{
			On: p1.ID,
			Quantity: quantity.Quantity{
				Quantity: 200.0,
				Unit:     "g",
			},
		},
	}

	p3 := NewComposition()
	p3.Unit = quantity.Quantity{500.0, "g"}
	p3.Dependencies = []Dependency{
		Dependency{
			On: p1.ID,
			Quantity: quantity.Quantity{
				Quantity: 100.0,
				Unit:     "g",
			},
		},
	}

	p4 := NewComposition()
	p4.Cost = 150.0
	p4.Unit = quantity.Quantity{100.0, "g"}

	p5 := NewComposition()
	p5.Unit = quantity.Quantity{1.0, "u"}
	p5.Dependencies = []Dependency{
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
	}

	p6 := NewComposition()
	p6.Unit = quantity.Quantity{2.0, "u"}
	p6.Dependencies = []Dependency{
		Dependency{
			On: p4.ID,
			Quantity: quantity.Quantity{
				Quantity: 350.0,
				Unit:     "g",
			},
		},
	}

	p7 := NewComposition()
	p7.Unit = quantity.Quantity{3.0, "u"}
	p7.Dependencies = []Dependency{
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
	r.Called(mock.Call("FindAll"))

	comps := make([]*Composition, 0)
	for _, c := range r.compositions {
		if c.Enabled {
			comps = append(comps, copyComposition(c))
		}
	}

	return comps, nil
}

func (r *mockRepository) FindByID(id string) (*Composition, error) {
	r.Called(mock.Call("FindByID", id))

	for _, c := range r.compositions {
		if c.ID.Hex() == id {
			return copyComposition(c), nil
		}
	}

	return nil, errors.NewInternal("NOT_FOUND").SetPath("composition/repository_mock.FindById")
}

func (r *mockRepository) FindUses(id string) ([]*Composition, error) {
	r.Called(mock.Call("FindUses", id))

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
	r.Called(mock.Call("FindByUsesUpdatedSinceLastChange", usesUpdated))

	comps := make([]*Composition, 0)
	for _, c := range r.compositions {
		if c.UsesUpdatedSinceLastChange == usesUpdated {
			comps = append(comps, copyComposition(c))
		}
	}

	return comps, nil
}

func (r *mockRepository) Insert(c *Composition) error {
	r.Called(mock.Call("Insert", c))

	c.UpdatedAt = time.Now()
	r.compositions = append(r.compositions, copyComposition(c))

	return nil
}

func (r *mockRepository) InsertMany(comps []*Composition) error {
	r.Called(mock.Call("InsertMany", comps))

	newComps := make([]*Composition, len(comps))
	for i, c := range comps {
		c.UpdatedAt = time.Now()
		newComps[i] = copyComposition(c)
	}
	r.compositions = append(r.compositions, newComps...)

	return nil
}

func (r *mockRepository) Update(c *Composition) error {
	r.Called(mock.Call("Update", c))

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
	r.Called(mock.Call("Delete", id))

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
