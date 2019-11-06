package composition

import (
	"errors"
)

type mockRepository struct {
	compositions []*Composition
}

func newMockRepository() *mockRepository {
	repo := &mockRepository{}
	return repo
}

func (r *mockRepository) FindAll() ([]*Composition, error) {
	comps := make([]*Composition, 0)
	for _, c := range r.compositions {
		if c.Enabled {
			comps = append(comps, copyComposition(c))
		}
	}
	return comps, nil
}

func (r *mockRepository) FindByID(id string) (*Composition, error) {
	for _, c := range r.compositions {
		if c.ID.String() == id && c.Enabled {
			return copyComposition(c), nil
		}
	}
	return nil, errors.New("Not found")
}

func (r *mockRepository) FindUses(id string) ([]*Composition, error) {
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
	return comps, nil
}

func (r *mockRepository) Insert(c *Composition) error {
	r.compositions = append(r.compositions, copyComposition(c))
	return nil
}

func (r *mockRepository) InsertMany(comps []*Composition) error {
	newComps := make([]*Composition, len(comps))
	for i, c := range comps {
		newComps[i] = copyComposition(c)
	}
	r.compositions = append(r.compositions, newComps...)
	return nil
}

func (r *mockRepository) Update(c *Composition) error {
	for _, comp := range r.compositions {
		if comp.ID.String() == c.ID.String() {
			if !comp.Enabled {
				return errors.New("Disabled")
			}
			*comp = *copyComposition(c)
			break
		}
	}
	return nil
}

func (r *mockRepository) Delete(id string) error {
	for _, comp := range r.compositions {
		if comp.ID.String() == id && comp.Enabled {
			comp.Enabled = false
			return nil
		}
	}
	return errors.New("Not found")
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

func (r *mockRepository) Clean() {
	r.compositions = make([]*Composition, 0)
}

func copyComposition(c *Composition) *Composition {
	copy := *c
	return &copy
}
