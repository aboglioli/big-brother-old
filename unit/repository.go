package unit

type Repository interface {
	FindAll() []*Unit
	FindByName(u string) *Unit
	FindByType(u string) []*Unit
	Exists(u string) bool
}

type repository struct {
	units map[string]*Unit
}

func NewRepository() Repository {
	return &repository{
		units: map[string]*Unit{
			"u": &Unit{"unit", "u", 1},

			"mg": &Unit{"mass", "mg", 0.001},
			"cg": &Unit{"mass", "cg", 0.01},
			"g":  &Unit{"mass", "g", 1},
			"kg": &Unit{"mass", "kg", 1000},

			"ml": &Unit{"volume", "ml", 0.001},
			"cl": &Unit{"volume", "cl", 0.01},
			"l":  &Unit{"volume", "l", 1},
			"kl": &Unit{"volume", "kl", 1000},

			"mm": &Unit{"length", "mm", 0.001},
			"cm": &Unit{"length", "cm", 0.01},
			"m":  &Unit{"length", "m", 1},
			"km": &Unit{"length", "km", 1000},
		},
	}
}

func (r *repository) FindAll() []*Unit {
	units := make([]*Unit, 0, len(r.units))
	for _, v := range r.units {
		units = append(units, v)
	}
	return units
}

func (r *repository) FindByName(n string) *Unit {
	return r.units[n]
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

func (r *repository) Exists(u string) bool {
	_, ok := r.units[u]
	return ok
}
