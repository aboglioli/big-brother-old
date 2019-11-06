package composition

import (
	"math"
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

func checkComp(t *testing.T, comps []*Composition, index int, shouldBe float64) {
	shouldBe = math.Round(shouldBe*100) / 100
	comp := comps[index]
	if comp.Cost != shouldBe {
		t.Errorf("Comp %d: %.2f should be %.2f", index, comp.Cost, shouldBe)
		for _, dep := range comp.Dependencies {
			t.Errorf("- dep %s subvalue %.2f", dep.Of.String(), dep.Subvalue)
		}
	}
}

func TestCreateComposition(t *testing.T) {
	repo := newMockRepository()
	serv := NewService(repo)

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

	t.Run("Non-existing dependency", func(t *testing.T) {
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

	t.Run("Incompatible dependency quantity", func(t *testing.T) {
		repo.Clean()
		dep, comp := newComposition(), newComposition()
		dep.Unit = quantity.Quantity{2, "kg"}
		repo.Insert(dep)
		comp.Dependencies = []*Dependency{
			&Dependency{
				Of: dep.ID,
				Quantity: quantity.Quantity{
					Quantity: 5,
					Unit:     "l",
				},
			},
		}

		if err := serv.Create(comp); !hasErrCode(err, "INCOMPATIBLE_DEPENDENCY_QUANTITY") {
			t.Error("Dependency cannot be created with incompatible dependency quantity")
		}
	})

	t.Run("Invalid dependency quantity", func(t *testing.T) {
		repo.Clean()
		dep, comp := newComposition(), newComposition()
		dep.Unit = quantity.Quantity{2, "kg"}
		repo.Insert(dep)
		comp.Dependencies = []*Dependency{
			&Dependency{
				Of: dep.ID,
				Quantity: quantity.Quantity{
					Quantity: 5,
					Unit:     "kk",
				},
			},
		}

		if err := serv.Create(comp); !hasErrCode(err, "INVALID_DEPENDENCY_QUANTITY") {
			t.Error("Dependency cannot be created with invalid dependency quantity")
		}

		comp.Dependencies[0].Quantity = quantity.Quantity{-5, "kg"}
		if err := serv.Create(comp); !hasErrCode(err, "INVALID_DEPENDENCY_QUANTITY") {
			t.Error("Dependency cannot be created with invalid dependency quantity")
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
			t.Error("Composition with single dependency should be created")
		}
	})

	t.Run("Calculate cost on creating", func(t *testing.T) {
		repo.Clean()
		dep, comp := newComposition(), newComposition()
		dep.Cost = 100
		dep.Unit = quantity.Quantity{
			Quantity: 2,
			Unit:     "kg",
		}
		comp.Dependencies = []*Dependency{
			&Dependency{
				Of: dep.ID,
				Quantity: quantity.Quantity{
					Quantity: 750,
					Unit:     "g",
				},
			},
		}
		repo.Insert(dep)

		if err := serv.Create(comp); err != nil {
			t.Error("Composition should be created")
		}

		if comp.Cost != 37.5 {
			t.Error("Cost wrong calculated")
		}
	})
}

func TestUpdateComposition(t *testing.T) {
	repo := newMockRepository()
	serv := NewService(repo)
	comps := makeMockedCompositions()
	repo.InsertMany(comps)
	for _, c := range comps {
		servImpl := serv.(*service)
		if err := servImpl.calculateDependenciesSubvalue(c.Dependencies); err != nil {
			t.Error(err)
			continue
		}
		c.CalculateCost()
		if err := repo.Update(c); err != nil {
			t.Error(err)
			continue
		}
	}

	t.Run("Update dependency", func(t *testing.T) {
		comps[0].Cost = 300
		comps[0].Unit = quantity.Quantity{
			Quantity: 2500,
			Unit:     "g",
		}

		if err := serv.Update(comps[0]); err != nil {
			t.Error(err)
		}

		comps, _ = repo.FindAll()
		if len(comps) != 7 {
			t.Error("Compositions count has changed")
		}

		c1 := 300.0
		q1 := 2.5
		c2 := 0.2 * c1 / q1 // 24
		c3 := 0.1 * c1 / q1 // 12
		c4 := 150.0
		c5 := 0.4*c2/0.2 + 0.05*c3/0.5 // 49.5
		c6 := 0.35 * c4 / 0.1          // 525
		c7 := 2*c5/1 + 1.5*c6/2        // 492.75

		checkComp(t, comps, 0, c1)
		checkComp(t, comps, 1, c2)
		checkComp(t, comps, 2, c3)
		checkComp(t, comps, 3, c4)
		checkComp(t, comps, 4, c5)
		checkComp(t, comps, 5, c6)
		checkComp(t, comps, 6, c7)
	})

	t.Run("Update dependency", func(t *testing.T) {
		comps[0].Cost = 300
		comps[0].Unit = quantity.Quantity{
			Quantity: 2500,
			Unit:     "g",
		}

		if err := serv.Update(comps[0]); err != nil {
			t.Error(err)
		}

		comps, _ = repo.FindAll()
		if len(comps) != 7 {
			t.Error("Compositions count has changed")
		}

		c1 := 300.0
		q1 := 2.5
		c2 := 0.2 * c1 / q1 // 24
		c3 := 0.1 * c1 / q1 // 12
		c4 := 150.0
		c5 := 0.4*c2/0.2 + 0.05*c3/0.5 // 49.5
		c6 := 0.35 * c4 / 0.1          // 525
		c7 := 2*c5/1 + 1.5*c6/2        // 492.75

		checkComp(t, comps, 0, c1)
		checkComp(t, comps, 1, c2)
		checkComp(t, comps, 2, c3)
		checkComp(t, comps, 3, c4)
		checkComp(t, comps, 4, c5)
		checkComp(t, comps, 5, c6)
		checkComp(t, comps, 6, c7)
	})
}

func TestDeleteComposition(t *testing.T) {
	repo := newMockRepository()
	serv := NewService(repo)
	comp, dep := newComposition(), newComposition()
	dep.Cost = 10
	dep.Unit = quantity.Quantity{1, "u"}
	repo.Insert(dep)
	comp.Dependencies = []*Dependency{
		&Dependency{
			Of:       dep.ID,
			Quantity: quantity.Quantity{2, "u"},
		},
	}
	repo.Insert(comp)

	if err := serv.Delete(dep.ID.String()); !hasErrCode(err, "COMPOSITION_USED_AS_DEPENDENCY") || repo.Count() != 2 {
		t.Error("Used dependency cannot be deleted")
	}

	if err := serv.Delete(comp.ID.String()); err != nil || repo.Count() != 1 {
		t.Error("Not used composition should be deleted")
	}
}

func TestCalculateDependenciesSubvalue(t *testing.T) {
	repo := newMockRepository()
	serv := NewService(repo)
	comps := makeMockedCompositions()
	repo.InsertMany(comps)
	for _, c := range comps {
		servImpl := serv.(*service)
		err := servImpl.calculateDependenciesSubvalue(c.Dependencies)
		if err != nil {
			t.Error(err)
			continue
		}
		c.CalculateCost()
	}

	c1 := 200.0
	q1 := 2.0
	c2 := 0.2 * c1 / q1 // 20
	c3 := 0.1 * c1 / q1 // 10
	c4 := 150.0
	c5 := 0.4*c2/0.2 + 0.05*c3/0.5 // 41
	c6 := 0.35 * c4 / 0.1          // 525
	c7 := 2*c5/1 + 1.5*c6/2        // 475.75

	checkComp(t, comps, 0, c1)
	checkComp(t, comps, 1, c2)
	checkComp(t, comps, 2, c3)
	checkComp(t, comps, 3, c4)
	checkComp(t, comps, 4, c5)
	checkComp(t, comps, 5, c6)
	checkComp(t, comps, 6, c7)
}
