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
		c.UpdateUses = false
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
	call := mock.Call("FindAll")

	comps := make([]*Composition, 0)
	for _, c := range r.compositions {
		if c.Enabled {
			comps = append(comps, copyComposition(c))
		}
	}

	r.Called(call.Return(comps, nil))
	return comps, nil
}

func (r *mockRepository) FindByID(id string) (*Composition, error) {
	call := mock.Call("FindByID", id)

	for _, c := range r.compositions {
		if c.ID.Hex() == id {
			r.Called(call.Return(copyComposition(c), nil))
			return copyComposition(c), nil
		}
	}

	err := errors.NewInternal("NOT_FOUND").SetPath("composition/mock.FindById")
	r.Called(call.Return(nil, err))
	return nil, err
}

func (r *mockRepository) FindUses(id string) ([]*Composition, error) {
	call := mock.Call("FindUses", id)

	comps := make([]*Composition, 0)
	for _, c := range r.compositions {
		for _, d := range c.Dependencies {
			if d.On.Hex() == id {
				comps = append(comps, copyComposition(c))
				break
			}
		}
	}

	r.Called(call.Return(comps, nil))
	return comps, nil
}

func (r *mockRepository) FindByUpdateUses(updateUses bool) ([]*Composition, error) {
	call := mock.Call("FindByUpdateUses", updateUses)

	comps := make([]*Composition, 0)
	for _, c := range r.compositions {
		if c.UpdateUses == updateUses {
			comps = append(comps, copyComposition(c))
		}
	}

	r.Called(call.Return(comps, nil))
	return comps, nil
}

func (r *mockRepository) Insert(c *Composition) error {
	call := mock.Call("Insert", c)

	c.UpdatedAt = time.Now()
	r.compositions = append(r.compositions, copyComposition(c))

	r.Called(call.Return(nil))
	return nil
}

func (r *mockRepository) InsertMany(comps []*Composition) error {
	call := mock.Call("InsertMany", comps)

	newComps := make([]*Composition, len(comps))
	for i, c := range comps {
		c.UpdatedAt = time.Now()
		newComps[i] = copyComposition(c)
	}
	r.compositions = append(r.compositions, newComps...)

	r.Called(call.Return(nil))
	return nil
}

func (r *mockRepository) Update(c *Composition) error {
	call := mock.Call("Update", c)

	for _, comp := range r.compositions {
		if comp.ID.Hex() == c.ID.Hex() {
			*comp = *copyComposition(c)
			comp.UpdatedAt = time.Now()
			break
		}
	}

	r.Called(call.Return(nil))
	return nil
}

func (r *mockRepository) SetUpdateUses(id string, updateUses bool) error {
	call := mock.Call("SetUpdateUses", id, updateUses)

	for _, comp := range r.compositions {
		if comp.ID.Hex() == id {
			comp.UpdatedAt = time.Now()
			comp.UpdateUses = updateUses
			break
		}
	}

	r.Called(call.Return(nil))
	return nil
}

func (r *mockRepository) Delete(id string) error {
	call := mock.Call("Delete", id)

	for _, comp := range r.compositions {
		if comp.ID.Hex() == id {
			comp.UpdatedAt = time.Now()
			comp.Enabled = false
			r.Called(call.Return(nil))
			return nil
		}
	}

	err := errors.NewInternal("NOT_FOUND").SetPath("composition/mock.Delete")
	r.Called(call.Return(err))
	return err
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
