package composition

import (
	"time"

	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/tests"
)

type mockRepository struct {
	tests.Mock
	compositions []*Composition
}

func newMockRepository() *mockRepository {
	return &mockRepository{}
}

// Helpers
func (r *mockRepository) Clean() {
	r.compositions = make([]*Composition, 0)
}

// Implementation
func (r *mockRepository) FindAll() ([]*Composition, errors.Error) {
	r.Called("FindAll")

	comps := make([]*Composition, 0)
	for _, c := range r.compositions {
		if c.Enabled {
			comps = append(comps, copyComposition(c))
		}
	}

	return comps, nil
}

func (r *mockRepository) FindByID(id string) (*Composition, errors.Error) {
	r.Called("FindByID", id)

	for _, c := range r.compositions {
		if c.ID.Hex() == id && c.Enabled {
			return copyComposition(c), nil
		}
	}

	return nil, errors.NewInternal().SetPath("composition/repository_mock.FindById").SetCode("NOT_FOUND")
}

func (r *mockRepository) FindUses(id string) ([]*Composition, errors.Error) {
	r.Called("FindUses", id)

	comps := make([]*Composition, 0)
	for _, c := range r.compositions {
		if !c.Enabled {
			break
		}

		for _, d := range c.Dependencies {
			if d.On.Hex() == id {
				comps = append(comps, copyComposition(c))
				break
			}
		}
	}

	return comps, nil
}

func (r *mockRepository) FindByUsesUpdatedSinceLastChange(usesUpdated bool) ([]*Composition, errors.Error) {
	r.Called("FindByUsesUpdatedSinceLastChange", usesUpdated)

	comps := make([]*Composition, 0)
	for _, c := range r.compositions {
		if !c.Enabled {
			break
		}

		if c.UsesUpdatedSinceLastChange == usesUpdated {
			comps = append(comps, copyComposition(c))
		}
	}

	return comps, nil
}

func (r *mockRepository) Insert(c *Composition) errors.Error {
	r.Called("Insert", c)

	c.UpdatedAt = time.Now()
	r.compositions = append(r.compositions, copyComposition(c))

	return nil
}

func (r *mockRepository) InsertMany(comps []*Composition) errors.Error {
	r.Called("InsertMany", comps)

	newComps := make([]*Composition, len(comps))
	for i, c := range comps {
		c.UpdatedAt = time.Now()
		newComps[i] = copyComposition(c)
	}
	r.compositions = append(r.compositions, newComps...)

	return nil
}

func (r *mockRepository) Update(c *Composition) errors.Error {
	r.Called("Update", c)

	for _, comp := range r.compositions {
		if comp.ID.Hex() == c.ID.Hex() {
			if !comp.Enabled {
				return errors.NewInternal().SetPath("composition/repository_mock.Update").SetCode("DISABLED")
			}
			*comp = *copyComposition(c)
			comp.UpdatedAt = time.Now()
			break
		}
	}

	return nil
}

func (r *mockRepository) Delete(id string) errors.Error {
	r.Called("Delete", id)

	for _, comp := range r.compositions {
		if comp.ID.Hex() == id && comp.Enabled {
			comp.UpdatedAt = time.Now()
			comp.Enabled = false
			return nil
		}
	}

	return errors.NewInternal().SetPath("composition/repository_mock.Delete").SetCode("NOT_FOUND")
}

func (r *mockRepository) Count() int {
	count := 0
	for _, c := range r.compositions {
		if c.Enabled {
			count++
		}
	}

	return count
}

func copyComposition(c *Composition) *Composition {
	copy := *c
	return &copy
}
