package unit

type Repository interface {
	FindByName(string) *Unit
	FindByType(string) []*Unit
}

type repository struct {
	units []*Unit
}

func NewRepository() Repository {
	return &repository{
		units: []*Unit{
			// Unit
			NewUnit("unit", "u", 1),
			// Mass
			NewUnit("mass", "mg", 0.001),
			NewUnit("mass", "cg", 0.01),
			NewUnit("mass", "g", 1),
			NewUnit("mass", "kg", 1000),
			// Volume
			NewUnit("volume", "ml", 0.001),
			NewUnit("volume", "cl", 0.01),
			NewUnit("volume", "l", 1),
			NewUnit("volume", "kl", 1000),
			// Length
			NewUnit("length", "mm", 0.001),
			NewUnit("length", "cm", 0.01),
			NewUnit("length", "m", 1),
			NewUnit("length", "km", 1000),
		},
	}
}

func (r *repository) FindByName(n string) *Unit {
	for _, v := range r.units {
		if v.Name == n {
			return v
		}
	}

	return nil
}

func (r *repository) FindByType(t string) []*Unit {
	units := make([]*Unit, 0)

	for _, v := range r.units {
		if v.Type == t {
			units = append(units, v)
		}
	}

	return units
}
