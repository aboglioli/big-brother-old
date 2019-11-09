package composition

import (
	"errors"
	"time"
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
func (r *mockRepository) FindAll() ([]*Composition, error) {
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

func (r *mockRepository) FindByID(id string) (*Composition, error) {
	r.sleep()

	for _, c := range r.compositions {
		if c.ID.String() == id && c.Enabled {
			return copyComposition(c), nil
		}
	}

	if r.CountRequests {
		r.Requests++
	}

	return nil, errors.New("Not found")
}

func (r *mockRepository) FindUses(id string) ([]*Composition, error) {
	r.sleep()

	comps := make([]*Composition, 0)
	for _, c := range r.compositions {
		if !c.Enabled {
			break
		}

		for _, d := range c.Dependencies {
			if d.Of.String() == id {
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

func (r *mockRepository) Insert(c *Composition) error {
	r.sleep()

	r.compositions = append(r.compositions, copyComposition(c))

	if r.CountRequests {
		r.Requests++
	}

	return nil
}

func (r *mockRepository) InsertMany(comps []*Composition) error {
	r.sleep()

	newComps := make([]*Composition, len(comps))
	for i, c := range comps {
		newComps[i] = copyComposition(c)
	}
	r.compositions = append(r.compositions, newComps...)

	if r.CountRequests {
		r.Requests++
	}

	return nil
}

func (r *mockRepository) Update(c *Composition) error {
	r.sleep()

	for _, comp := range r.compositions {
		if comp.ID.String() == c.ID.String() {
			if !comp.Enabled {
				return errors.New("Disabled")
			}
			*comp = *copyComposition(c)
			break
		}
	}

	if r.CountRequests {
		r.Requests++
	}

	return nil
}

func (r *mockRepository) Delete(id string) error {
	r.sleep()

	for _, comp := range r.compositions {
		if comp.ID.String() == id && comp.Enabled {
			comp.Enabled = false
			return nil
		}
	}

	if r.CountRequests {
		r.Requests++
	}

	return errors.New("Not found")
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
