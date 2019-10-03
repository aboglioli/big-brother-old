package main

import (
	"context"

	"github.com/aboglioli/big-brother/infrastructure/db"
)

type Composition struct {
	ID           string       `bson:"_id" validate:"required"`
	Cost         float64      `bson:"cost" validate:"required"`
	Quantity     Quantity     `bson:"quantity" validate:"required"`
	Dependencies []Dependency `bson:"dependencies" validate:"required"`
}

type Quantity struct {
	Value float64 `bson:"value" validate:"required"`
	Unit  string  `bson:"unit" validate:"required"`
}

type Dependency struct {
	From     string   `bson:"from" validate:"required"`
	Quantity Quantity `bson:"quantity" validate:"required"`
	Subvalue float64  `bson:"subvalue" validate:"required"`
}

func main() {
	d, _ := db.Get("Product")
	coll := d.Collection("Composition")
	coll.Drop(context.TODO())

	compositions := []Composition{
		Composition{
			ID:   "P1",
			Cost: 200.0,
			Quantity: Quantity{
				Value: 2.0,
				Unit:  "kg",
			},
		},
		Composition{
			ID: "P2",
			Quantity: Quantity{
				Value: 0.2,
				Unit:  "kg",
			},
			Dependencies: []Dependency{
				Dependency{
					From: "P1",
					Quantity: Quantity{
						Value: 200.0,
						Unit:  "g",
					},
				},
			},
		},
		Composition{
			ID: "P3",
			Quantity: Quantity{
				Value: 500.0,
				Unit:  "g",
			},
			Dependencies: []Dependency{
				Dependency{
					From: "P1",
					Quantity: Quantity{
						Value: 100.0,
						Unit:  "g",
					},
				},
				Dependency{
					From: "P10",
					Quantity: Quantity{
						Value: 0.1,
						Unit:  "l",
					},
				},
			},
		},
		Composition{
			ID:   "P4",
			Cost: 150.0,
			Quantity: Quantity{
				Value: 100.0,
				Unit:  "g",
			},
		},
		Composition{
			ID: "P5",
			Quantity: Quantity{
				Value: 1.0,
				Unit:  "u",
			},
			Dependencies: []Dependency{
				Dependency{
					From: "P2",
					Quantity: Quantity{
						Value: 400.0,
						Unit:  "g",
					},
				},
				Dependency{
					From: "P3",
					Quantity: Quantity{
						Value: 50.0,
						Unit:  "g",
					},
				},
			},
		},
		Composition{
			ID: "P6",
			Quantity: Quantity{
				Value: 2.0,
				Unit:  "u",
			},
			Dependencies: []Dependency{
				Dependency{
					From: "P4",
					Quantity: Quantity{
						Value: 20.0,
						Unit:  "g",
					},
				},
			},
		},
		Composition{
			ID: "P7",
			Quantity: Quantity{
				Value: 3.0,
				Unit:  "u",
			},
			Dependencies: []Dependency{
				Dependency{
					From: "P5",
					Quantity: Quantity{
						Value: 2.0,
						Unit:  "u",
					},
				},
				Dependency{
					From: "P6",
					Quantity: Quantity{
						Value: 1.5,
						Unit:  "u",
					},
				},
			},
		},
		Composition{
			ID: "P8",
			Quantity: Quantity{
				Value: 1.0,
				Unit:  "u",
			},
			Dependencies: []Dependency{
				Dependency{
					From: "P7",
					Quantity: Quantity{
						Value: 2.5,
						Unit:  "u",
					},
				},
			},
		},
		Composition{
			ID: "P9",
			Quantity: Quantity{
				Value: 3,
				Unit:  "u",
			},
			Dependencies: []Dependency{
				Dependency{
					From: "P3",
					Quantity: Quantity{
						Value: 0.25,
						Unit:  "kg",
					},
				},
				Dependency{
					From: "P10",
					Quantity: Quantity{
						Value: 1,
						Unit:  "l",
					},
				},
			},
		},
		Composition{
			ID: "P10",
			Quantity: Quantity{
				Value: 500,
				Unit:  "ml",
			},
		},
	}

	rawComps := make([]interface{}, len(compositions))
	for i, c := range compositions {
		rawComps[i] = c
	}
	coll.InsertMany(context.TODO(), rawComps)
}
