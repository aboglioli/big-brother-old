package composition

import (
	"time"

	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/tests/mock"
)

type mockRepository struct {
	mock.Mock
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
