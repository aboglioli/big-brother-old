package unit

type Unit struct {
	Type     string
	Name     string
	Modifier float64
}

func NewUnit(t string, n string, m float64) *Unit {
	return &Unit{
		Type:     t,
		Name:     n,
		Modifier: m,
	}
}

func (u1 *Unit) Equals(u2 *Unit) bool {
	return u1.Type == u2.Type && u1.Name == u2.Name && u1.Modifier == u2.Modifier
}

func (u1 *Unit) SameType(u2 *Unit) bool {
	return u1.Type == u2.Type
}
