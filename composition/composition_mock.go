package composition

import (
	"github.com/aboglioli/big-brother/quantity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func newComposition() *Composition {
	comp := NewComposition()
	comp.Unit.Unit = "u"
	comp.Stock.Unit = "u"
	return comp
}

func makeCompositions() []*Composition {
	p1 := &Composition{
		ID:   primitive.NewObjectID(),
		Cost: 200.0,
		Unit: quantity.Quantity{
			Quantity: 2.0,
			Unit:     "kg",
		},
	}
	p2 := &Composition{
		ID:   primitive.NewObjectID(),
		Cost: 0.0, // 0.2 * 200 / 2 = 20
		Unit: quantity.Quantity{
			Quantity: 0.2,
			Unit:     "kg",
		},
		Dependencies: []*Dependency{
			&Dependency{
				Of: p1.ID,
				Quantity: quantity.Quantity{
					Quantity: 200.0,
					Unit:     "g",
				},
			},
		},
	}
	p3 := &Composition{
		ID:   primitive.NewObjectID(),
		Cost: 0.0, // 0.1 * 200 / 2 = 10
		Unit: quantity.Quantity{
			Quantity: 500.0,
			Unit:     "g",
		},
		Dependencies: []*Dependency{
			&Dependency{
				Of: p1.ID,
				Quantity: quantity.Quantity{
					Quantity: 100.0,
					Unit:     "g",
				},
			},
		},
	}
	p4 := &Composition{
		ID:   primitive.NewObjectID(),
		Cost: 150.0,
		Unit: quantity.Quantity{
			Quantity: 100.0,
			Unit:     "g",
		},
	}
	p5 := &Composition{
		ID:   primitive.NewObjectID(),
		Cost: 0.0, // 0.4*20/0.2 + 0.05*10/0.5 = 41
		Unit: quantity.Quantity{
			Quantity: 1.0,
			Unit:     "u",
		},
		Dependencies: []*Dependency{
			&Dependency{
				Of: p2.ID,
				Quantity: quantity.Quantity{
					Quantity: 400.0,
					Unit:     "g",
				},
			},
			&Dependency{
				Of: p3.ID,
				Quantity: quantity.Quantity{
					Quantity: 50.0,
					Unit:     "g",
				},
			},
		},
	}
	p6 := &Composition{
		ID:   primitive.NewObjectID(),
		Cost: 0.0, // 0.35*150/0.1 = 525
		Unit: quantity.Quantity{
			Quantity: 2.0,
			Unit:     "u",
		},
		Dependencies: []*Dependency{
			&Dependency{
				Of: p4.ID,
				Quantity: quantity.Quantity{
					Quantity: 350.0,
					Unit:     "g",
				},
			},
		},
	}
	p7 := &Composition{
		ID:   primitive.NewObjectID(),
		Cost: 0.0, // 2*41/1 + 1.5*525/2 = 475.75
		Unit: quantity.Quantity{
			Quantity: 3.0,
			Unit:     "u",
		},
		Dependencies: []*Dependency{
			&Dependency{
				Of: p5.ID,
				Quantity: quantity.Quantity{
					Quantity: 2.0,
					Unit:     "u",
				},
			},
			&Dependency{
				Of: p6.ID,
				Quantity: quantity.Quantity{
					Quantity: 1.5,
					Unit:     "u",
				},
			},
		},
	}
	comps := []*Composition{p1, p2, p3, p4, p5, p6, p7}

	for _, c := range comps {
		c.AutoupdateCost = true
		c.Enabled = true
	}

	return comps
}
