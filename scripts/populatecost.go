package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"

	"github.com/aboglioli/big-brother/infrastructure/db"
	"github.com/aboglioli/big-brother/quantity"
	"github.com/aboglioli/big-brother/unit"
	"go.mongodb.org/mongo-driver/mongo"
)

type Dependency struct {
	ID       string            `bson:"id" validate:"required"`
	Quantity quantity.Quantity `bson:"quantity" validate:"required"`
	Subvalue float64           `bson:"subvalue"`
}

type Cost struct {
	ID           string            `bson:"_id" validate:"required"`
	Cost         float64           `bson:"cost" validate:"required"`
	Quantity     quantity.Quantity `bson:"quantity" validate:"required"`
	Dependencies []*Dependency     `bson:"dependencies" validate:"required"`
}

var unitRepo unit.Repository
var coll *mongo.Collection
var id int

func init() {
	unitRepo = unit.NewRepository()
	qServ = quantity.NewService()
	d, _ := db.Get("product")
	coll = d.Collection("cost")
}

func main() {
	coll.Drop(context.Background())

	// Inputs
	var inputs []*Cost
	for i := 0; i < 5; i++ {
		inputs = append(inputs, &Cost{
			ID:       getID(),
			Cost:     randCostValue(),
			Quantity: randomQuantity(),
		})
	}

	// Simple dependency
	var simple []*Cost
	for i := 0; i < 5; i++ {
		dep := inputs[rand.Intn(len(inputs))]
		cost := &Cost{
			ID:       getID(),
			Cost:     0,
			Quantity: randomQuantity(),
			Dependencies: []*Dependency{
				&Dependency{
					ID:       dep.ID,
					Quantity: randomQuantityByUnitType(dep.Quantity.Unit),
				},
			},
		}

		simple = append(simple, cost)
	}

	// Complex dependency
	var complex []*Cost
	for i := 0; i < 5; i++ {
		cost := &Cost{
			ID:       getID(),
			Cost:     0,
			Quantity: randomQuantity(),
		}

		for j := 0; j < rand.Intn(5)+1; j++ {
			var dep *Cost
			r := rand.Int()
			if r%2 == 0 {
				dep = inputs[rand.Intn(len(inputs))]
			} else {
				dep = simple[rand.Intn(len(simple))]
			}

			cost.Dependencies = append(cost.Dependencies, &Dependency{
				ID:       dep.ID,
				Quantity: randomQuantityByUnitType(dep.Quantity.Unit),
			})
		}

		complex = append(complex, cost)
	}

	// Insert costs
	costs := append(inputs, simple...)
	costs = append(costs, complex...)
	rawCosts := make([]interface{}, len(costs))
	for i, c := range costs {
		rawCosts[i] = c
	}
	r, err := coll.InsertMany(context.Background(), rawCosts)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(r.InsertedIDs))
}

func getID() string {
	c := id
	id++
	return fmt.Sprintf("Product-%d", c)
}

func randFloat(max float64, dec int) float64 {
	d := math.Pow10(dec)
	v := rand.Float64() * max

	return math.Round(v*d) / d
}

func randIntAsFloat(max int, m int) float64 {
	v := rand.Intn(max)

	return float64(v - (v % m))
}

func randCostValue() float64 {
	// return randFloat(100, 0)
	return randIntAsFloat(1000, 5)
}

func randomUnit() *unit.Unit {
	units := unitRepo.FindAll()
	return units[rand.Intn(len(units))]
}

func randomUnitByType(u string) *unit.Unit {
	unit := unitRepo.FindByName(u)
	units := unitRepo.FindByType(unit.Type)
	return units[rand.Intn(len(units))]
}

func randomQuantity() quantity.Quantity {
	// value := randFloat(1000, 0)
	value := randIntAsFloat(1000, 5) + 5
	unit := randomUnit()

	return quantity.Quantity{
		Unit:     unit.Name,
		Quantity: value,
	}
}

func randomQuantityByUnitType(u string) quantity.Quantity {
	// value := randFloat(1000, 0)
	value := randIntAsFloat(1000, 5) + 5
	unit := randomUnitByType(u)

	return quantity.Quantity{
		Unit:     unit.Name,
		Quantity: value,
	}
}
