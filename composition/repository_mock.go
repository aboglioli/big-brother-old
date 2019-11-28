package composition

import (
	"time"

	"github.com/aboglioli/big-brother/pkg/errors"
)

type mockRepository struct {
	compositions  []*Composition
	CountRequests bool
	Requests      int
	SimulateDelay bool
	SleepTime     int
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		CountRequests: false,
		Requests:      0,
		SimulateDelay: false,
		SleepTime:     0,
	}
}

// Helpers
func (r *mockRepository) Clean() {
	r.compositions = make([]*Composition, 0)
}

func (r *mockRepository) sleep() {
	if r.SimulateDelay {
		time.Sleep(time.Duration(r.SleepTime) * time.Millisecond)
	}
}

// Implementation
func (r *mockRepository) FindAll() ([]*Composition, errors.Error) {
	r.sleep()

	comps := make([]*Composition, 0)
	for _, c := range r.compositions {
		if c.Enabled {
			comps = append(comps, copyComposition(c))
		}
	}

	if r.CountRequests {
		r.Requests++
	}

	return comps, nil
}

func (r *mockRepository) FindByID(id string) (*Composition, errors.Error) {
	r.sleep()

	for _, c := range r.compositions {
		if c.ID.Hex() == id && c.Enabled {
			return copyComposition(c), nil
		}
	}

	if r.CountRequests {
		r.Requests++
	}

	return nil, errors.NewInternal().SetPath("composition/repository_mock.FindById").SetCode("NOT_FOUND")
}

func (r *mockRepository) FindUses(id string) ([]*Composition, errors.Error) {
	r.sleep()

	comps := make([]*Composition, 0)
	for _, c := range r.compositions {
		if !c.Enabled {
			break
		}

		for _, d := range c.Dependencies {
			if d.Of.Hex() == id {
				comps = append(comps, copyComposition(c))
				break
			}
		}
	}

	if r.CountRequests {
		r.Requests++
	}

	return comps, nil
}

func (r *mockRepository) FindByUsesUpdatedSinceLastChange(usesUpdated bool) ([]*Composition, errors.Error) {
	r.sleep()

	comps := make([]*Composition, 0)
	for _, c := range r.compositions {
		if !c.Enabled {
			break
		}

		if c.UsesUpdatedSinceLastChange == usesUpdated {
			comps = append(comps, copyComposition(c))
		}
	}

	if r.CountRequests {
		r.Requests++
	}

	return comps, nil
}

func (r *mockRepository) Insert(c *Composition) errors.Error {
	r.sleep()

	c.UpdatedAt = time.Now()
	r.compositions = append(r.compositions, copyComposition(c))

	if r.CountRequests {
		r.Requests++
	}

	return nil
}

func (r *mockRepository) InsertMany(comps []*Composition) errors.Error {
	r.sleep()

	newComps := make([]*Composition, len(comps))
	for i, c := range comps {
		c.UpdatedAt = time.Now()
		newComps[i] = copyComposition(c)
	}
	r.compositions = append(r.compositions, newComps...)

	if r.CountRequests {
		r.Requests++
	}

	return nil
}

func (r *mockRepository) Update(c *Composition) errors.Error {
	r.sleep()

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

	if r.CountRequests {
		r.Requests++
	}

	return nil
}

func (r *mockRepository) Delete(id string) errors.Error {
	r.sleep()

	for _, comp := range r.compositions {
		if comp.ID.Hex() == id && comp.Enabled {
			comp.UpdatedAt = time.Now()
			comp.Enabled = false
			return nil
		}
	}

	if r.CountRequests {
		r.Requests++
	}

	return errors.NewInternal().SetPath("composition/repository_mock.Delete").SetCode("NOT_FOUND")
}

func (r *mockRepository) Count() int {
	r.sleep()

	count := 0
	for _, c := range r.compositions {
		if c.Enabled {
			count++
		}
	}

	if r.CountRequests {
		r.Requests++
	}

	return count
}

func copyComposition(c *Composition) *Composition {
	copy := *c
	return &copy
}
