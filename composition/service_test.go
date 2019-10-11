package composition

import (
	"testing"

	"github.com/aboglioli/big-brother/quantity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var comps []*Composition

func init() {
	comps = []*Composition{}
}

func newComposition() *Composition {
	comp := NewComposition()
	comp.Unit.Unit = "u"
	comp.Stock.Unit = "u"
	return comp
}

func TestCreate(t *testing.T) {
	repo := NewMockRepository()
	qServ := quantity.NewService()
	serv := NewService(repo, qServ)

	t.Run("Default values with valid units", func(t *testing.T) {
		comp := newComposition()
		err := serv.Create(comp)

		if err != nil || repo.Count() != 1 {
			t.Error("Composition should be created")
		}
	})

	t.Run("Negative cost", func(t *testing.T) {
		comp := newComposition()
		comp.Cost = -1.0
		err := serv.Create(comp)

		if err == nil {
			t.Error("Cost can't be negative")
		}
	})

	t.Run("Invalid units", func(t *testing.T) {
		comp := newComposition()
		comp.Unit.Unit = ""
		err := serv.Create(comp)
		if err == nil {
			t.Error("Unit shuld exist")
		}

		comp.Unit.Unit = "u"
		comp.Stock.Unit = ""
		err = serv.Create(comp)
		if err == nil {
			t.Error("Stock unit should exist")
		}
	})

	t.Run("Invalid dependency", func(t *testing.T) {
		comp := newComposition()
		comp.Dependencies = append(comp.Dependencies, &Dependency{
			Of: primitive.NewObjectID(),
			Quantity: quantity.Quantity{
				Quantity: 5,
				Unit:     "u",
			},
		})
		err := serv.Create(comp)

		if err == nil {
			t.Error("Dependency doesn't exist")
		}
	})
}
