package composition

import (
	"github.com/aboglioli/big-brother/quantity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MakeCompositions() []*Composition {
	p1 := &Composition{
		ID:   primitive.NewObjectID(),
		Cost: 200.0,
		Unit: quantity.Quantity{
			Quantity: 2.0,
			Unit:     "kg",
		},
	}
	p2 := &Composition{
		ID: primitive.NewObjectID(),
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
		ID: primitive.NewObjectID(),
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
		ID: primitive.NewObjectID(),
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
		ID: primitive.NewObjectID(),
		Unit: quantity.Quantity{
			Quantity: 2.0,
			Unit:     "u",
		},
		Dependencies: []*Dependency{
			&Dependency{
				Of: p4.ID,
				Quantity: quantity.Quantity{
					Quantity: 20.0,
					Unit:     "g",
				},
			},
		},
	}
	p7 := &Composition{
		ID: primitive.NewObjectID(),
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
	return []*Composition{p1, p2, p3, p4, p5, p6, p7}
}
