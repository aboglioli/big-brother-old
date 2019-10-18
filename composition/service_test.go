package composition

import (
	"testing"

	"github.com/aboglioli/big-brother/errors"
	"github.com/aboglioli/big-brother/quantity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func hasErrCode(err errors.Error, code string) bool {
	if err == nil {
		return false
	}
	return err.Code() == code
}

func TestCreate(t *testing.T) {
	repo := NewMockRepository()
	qServ := quantity.NewService()
	serv := NewService(repo, qServ)

	// Errors
	t.Run("Negative cost", func(t *testing.T) {
		comp := newComposition()
		comp.Cost = -1.0
		err := serv.Create(comp)

		if !hasErrCode(err, "NEGATIVE_COST") {
			t.Error("Cost can't be negative")
		}
	})

	t.Run("Invalid units", func(t *testing.T) {
		comp := newComposition()
		comp.Unit.Unit = ""
		err := serv.Create(comp)
		if !hasErrCode(err, "INVALID_UNIT") {
			t.Error("Unit shuld exist")
		}

		comp.Unit.Unit = "u"
		comp.Stock.Unit = ""
		err = serv.Create(comp)
		if !hasErrCode(err, "INVALID_STOCK") {
			t.Error("Stock unit should exist")
		}
	})

	t.Run("Invalid dependency", func(t *testing.T) {
		comp := newComposition()
		comp.Dependencies = []*Dependency{
			&Dependency{
				Of: primitive.NewObjectID(),
				Quantity: quantity.Quantity{
					Quantity: 5,
					Unit:     "u",
				},
			},
		}
		err := serv.Create(comp)

		if !hasErrCode(err, "DEPENDENCY_DOES_NOT_EXIST") {
			t.Error("Check dependency existence")
		}
	})

	// Create
	t.Run("Default values with valid units", func(t *testing.T) {
		repo.Clean()
		comp := newComposition()
		err := serv.Create(comp)

		if err != nil || repo.Count() != 1 {
			t.Error("Composition should be created")
		}
	})

	t.Run("Valid dependency", func(t *testing.T) {
		repo.Clean()
		dep, comp := newComposition(), newComposition()
		repo.Insert(dep)
		comp.Dependencies = []*Dependency{
			&Dependency{
				Of: dep.ID,
				Quantity: quantity.Quantity{
					Quantity: 5,
					Unit:     "u",
				},
			},
		}

		err := serv.Create(comp)

		if err != nil || repo.Count() != 2 {
			t.Error("Component with single dependency should be created")
		}
	})
}
