package composition

import (
	"math"
	"testing"

	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/events"
	"github.com/aboglioli/big-brother/pkg/quantity"
	"github.com/aboglioli/big-brother/pkg/tests"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func hasErrCode(err errors.Error, code string) bool {
	if err == nil {
		return false
	}
	return err.Code() == code
}

func checkCompCost(t *testing.T, comps []*Composition, index int, costShouldBe float64) {
	costShouldBe = math.Round(costShouldBe*1000) / 1000
	comp := comps[index]
	if comp.Cost != costShouldBe {
		t.Errorf("Comp %d: %.2f should be %.2f", index, comp.Cost, costShouldBe)
		for _, dep := range comp.Dependencies {
			t.Errorf("- dep %s subvalue %.2f", dep.On.Hex(), dep.Subvalue)
		}
	}
}

func compToCreateRequest(c *Composition) *CreateRequest {
	id := c.ID.Hex()
	return &CreateRequest{
		ID:             &id,
		Name:           c.Name,
		Cost:           c.Cost,
		Unit:           c.Unit,
		Stock:          &c.Stock,
		Dependencies:   c.Dependencies,
		AutoupdateCost: &c.AutoupdateCost,
	}
}

func compToUpdateRequest(c *Composition) *UpdateRequest {
	id := c.ID.Hex()
	return &UpdateRequest{
		ID:             &id,
		Name:           &c.Name,
		Cost:           &c.Cost,
		Unit:           &c.Unit,
		Stock:          &c.Stock,
		Dependencies:   c.Dependencies,
		AutoupdateCost: &c.AutoupdateCost,
	}
}

func TestCreateComposition(t *testing.T) {
	repo, eventMgr := newMockRepository(), events.GetMockManager()
	serv := NewService(repo, eventMgr)

	t.Run("Non-existing dependency", func(t *testing.T) {
		comp := newComposition()
		comp.Dependencies = []Dependency{
			Dependency{
				On: primitive.NewObjectID(),
				Quantity: quantity.Quantity{
					Quantity: 5,
					Unit:     "u",
				},
			},
		}
		_, err := serv.Create(compToCreateRequest(comp))

		tests.ErrCode(t, err, "DEPENDENCY_DOES_NOT_EXIST", "Check dependency existence")
	})

	t.Run("Incompatible dependency quantity", func(t *testing.T) {
		repo.Clean()
		dep, comp := newComposition(), newComposition()
		dep.Unit = quantity.Quantity{2, "kg"}
		repo.Insert(dep)
		comp.Dependencies = []Dependency{
			Dependency{
				On: dep.ID,
				Quantity: quantity.Quantity{
					Quantity: 5,
					Unit:     "l",
				},
			},
		}

		_, err := serv.Create(compToCreateRequest(comp))
		tests.ErrCode(t, err, "INCOMPATIBLE_DEPENDENCY_QUANTITY", "Dependency cannot be created with incompatible dependency quantity")
	})

	// Create
	t.Run("Default values with valid units and raise event 'CompositionCreated'", func(t *testing.T) {
		repo.Clean()
		eventMgr.Clean()
		comp := newComposition()
		_, err := serv.Create(compToCreateRequest(comp))

		tests.Ok(t, err, "Should be created")
		tests.Equal(t, repo.Count(), 1, "Should be created")
		tests.Equal(t, eventMgr.Count(), 1, "Should emit an event")

		msgs := eventMgr.Messages()
		msg := msgs[0]
		tests.Equal(t, msg.Type(), "CompositionCreated", "Wrong event")
	})

	t.Run("Assign stock automatically from unit", func(t *testing.T) {
		repo.Clean()
		comp := newComposition()
		comp.Unit = quantity.Quantity{2, "kg"}
		createReq := compToCreateRequest(comp)
		createReq.Stock = nil
		c, err := serv.Create(createReq)

		tests.Ok(t, err, "No error")
		tests.Assert(t, c.Stock.Equals(quantity.Quantity{0, c.Unit.Unit}), "Stock should be auto-assigned")
	})

	t.Run("Valid dependency", func(t *testing.T) {
		repo.Clean()
		dep, comp := newComposition(), newComposition()
		repo.Insert(dep)
		comp.Dependencies = []Dependency{
			Dependency{
				On: dep.ID,
				Quantity: quantity.Quantity{
					Quantity: 5,
					Unit:     "u",
				},
			},
		}

		_, err := serv.Create(compToCreateRequest(comp))

		tests.Ok(t, err, "No error")
		tests.Equal(t, repo.Count(), 2, "Composition with single dependency should be created")
	})

	t.Run("Calculate cost on creating and comparte with raised event", func(t *testing.T) {
		repo.Clean()
		eventMgr.Clean()
		dep, comp := newComposition(), newComposition()
		dep.Cost = 100
		dep.Unit = quantity.Quantity{
			Quantity: 2,
			Unit:     "kg",
		}
		comp.Dependencies = []Dependency{
			Dependency{
				On: dep.ID,
				Quantity: quantity.Quantity{
					Quantity: 750,
					Unit:     "g",
				},
			},
		}
		repo.Insert(dep)

		c, err := serv.Create(compToCreateRequest(comp))

		tests.Ok(t, err, "No error")
		tests.Equal(t, c.Cost, 37.5, "Cost not calculated")
		tests.Equal(t, eventMgr.Count(), 1, "Should raise an event")

		msgs := eventMgr.Messages()
		msg := msgs[0]
		tests.Assert(t, msg.Type() == "CompositionCreated" && msg.Key == "composition.created", "Wrong event")

		var evt CompositionChangedEvent
		tests.Ok(t, msg.Decode(&evt), "Decode")
		tests.Assert(t, evt.Composition.Cost == c.Cost && evt.Composition.ID.Hex() == c.ID.Hex(), "Composition from event is not the expected one")
	})
}

func TestUpdateComposition(t *testing.T) {
	repo, eventMgr := newMockRepository(), events.GetMockManager()
	serv := NewService(repo, eventMgr)

	t.Run("Update dependency and raise events", func(t *testing.T) {
		repo.Clean()
		eventMgr.Clean()

		comps := makeMockedCompositions()
		repo.InsertMany(comps)
		for _, c := range comps {
			servImpl := serv.(*service)
			deps, err := servImpl.calculateDependenciesSubvalues(c.Dependencies)
			tests.Ok(t, err, "No error")
			c.SetDependencies(deps)
			tests.Ok(t, repo.Update(c), "No error")
		}

		c := comps[0]
		c.Cost = 300
		c.Unit = quantity.Quantity{
			Quantity: 2500,
			Unit:     "g",
		}

		c, err := serv.Update(c.ID.Hex(), compToUpdateRequest(c))
		tests.Ok(t, err, "No error")
		updatedUses, err := serv.UpdateUses(c)
		tests.Ok(t, err, "No error")

		tests.Equal(t, len(updatedUses), 4, "Uses weren't updated")

		comps, _ = repo.FindAll()
		tests.Equal(t, len(comps), 7, "Compositions count has changed")

		c1 := 300.0
		q1 := 2.5
		c2 := 0.2 * c1 / q1 // 24
		c3 := 0.1 * c1 / q1 // 12
		c4 := 150.0
		c5 := 0.4*c2/0.2 + 0.05*c3/0.5 // 49.2
		c6 := 0.35 * c4 / 0.1          // 525
		c7 := 2*c5/1 + 1.5*c6/2        // 492.75

		checkCompCost(t, comps, 0, c1)
		checkCompCost(t, comps, 1, c2)
		checkCompCost(t, comps, 2, c3)
		checkCompCost(t, comps, 3, c4)
		checkCompCost(t, comps, 4, c5)
		checkCompCost(t, comps, 5, c6)
		checkCompCost(t, comps, 6, c7)

		// Check events
		tests.Equal(t, eventMgr.Count(), 2, "Should raise events")

		msgs := eventMgr.Messages()
		tests.Assert(t, msgs[0].Type() == "CompositionUpdatedManually" && msgs[0].Key == "composition.updated", "Wrong event")

		var compUpdatedManuallyEvent CompositionChangedEvent
		tests.Ok(t, msgs[0].Decode(&compUpdatedManuallyEvent), "Error")
		tests.Equal(t, compUpdatedManuallyEvent.Composition.ID.Hex(), c.ID.Hex(), "Different composition")
		tests.Assert(t, msgs[1].Type() == "CompositionsUpdatedAutomatically" && msgs[0].Key == "composition.updated", "Wrong event")

		var compsUpdatedAutomaticallyEvent CompositionsUpdatedAutomaticallyEvent
		tests.Ok(t, msgs[1].Decode(&compsUpdatedAutomaticallyEvent), "No error")
		tests.Equal(t, len(compsUpdatedAutomaticallyEvent.Compositions), 4, "Update automatically")
	})

	t.Run("Creating and updating", func(t *testing.T) {
		repo.Clean()
		comp := newComposition()
		comp.Cost = 30
		comp.Unit = quantity.Quantity{1, "u"}
		comp.Stock = quantity.Quantity{1, "u"}

		createdComp, err := serv.Create(compToCreateRequest(comp))
		tests.Ok(t, err, "No Error")

		updatedComp, err := serv.Update(createdComp.ID.Hex(), compToUpdateRequest(createdComp))
		tests.Ok(t, err, "No error")
		tests.Equal(t, createdComp.ID.Hex(), updatedComp.ID.Hex(), "ID changed")
	})

	t.Run("Invalid units", func(t *testing.T) {
		repo.Clean()
		comp := newComposition()
		comp.Cost = 30
		comp.Unit = quantity.Quantity{1, "u"}
		comp.Stock = quantity.Quantity{1, "u"}

		createdComp, err := serv.Create(compToCreateRequest(comp))
		tests.Ok(t, err, "No error")

		createdComp.Unit = quantity.Quantity{1, "asd"}
		_, err = serv.Update(createdComp.ID.Hex(), compToUpdateRequest(createdComp))
		tests.ErrCode(t, err, "INVALID_UNIT", "Should return error due to invalid unit")

		createdComp.Unit = quantity.Quantity{1, "u"}
		createdComp.Stock = quantity.Quantity{1, "asd"}
		_, err = serv.Update(createdComp.ID.Hex(), compToUpdateRequest(createdComp))
		tests.ErrCode(t, err, "INVALID_STOCK", "Should return error due to invalid unit in stock")

		createdComp.Unit = quantity.Quantity{1, "kg"}
		createdComp.Stock = quantity.Quantity{1, "l"}
		_, err = serv.Update(createdComp.ID.Hex(), compToUpdateRequest(createdComp))
		tests.ErrCode(t, err, "INCOMPATIBLE_STOCK_AND_UNIT", "Stock and unit should be compatible")
	})

	t.Run("Empty stock ignored on updating", func(t *testing.T) {
		repo.Clean()
		comp := newComposition()
		comp.Unit = quantity.Quantity{5, "l"}
		comp.Stock = quantity.Quantity{25, "l"}
		repo.Insert(comp)

		updateReq := compToUpdateRequest(comp)
		updateReq.Stock = nil

		c, err := serv.Update(comp.ID.Hex(), updateReq)
		tests.Ok(t, err, "No error")

		tests.Assert(t, c.Stock.Equals(quantity.Quantity{25, "l"}), "Empty stock should be ignored")

		updateReq.Stock = &quantity.Quantity{4000, "ml"}
		c, err = serv.Update(comp.ID.Hex(), updateReq)
		tests.Ok(t, err, "No error")

		tests.Assert(t, c.Stock.Equals(quantity.Quantity{4000, "ml"}), "Non-empty stock should be assigned")
	})
}

func TestCreateAndUpdateDependencies(t *testing.T) {
	repo, eventMgr := newMockRepository(), events.GetMockManager()
	serv := NewService(repo, eventMgr)

	repo.Clean()
	comp, dep1, dep2, dep3 := newComposition(), newComposition(), newComposition(), newComposition()
	dep1.Cost = 100
	dep1.Unit = quantity.Quantity{1, "kg"}
	dep2.Cost = 200
	dep2.Unit = quantity.Quantity{4000, "g"}
	dep3.Cost = 75
	dep3.Unit = quantity.Quantity{0.6, "kg"}
	comp.Dependencies = []Dependency{
		Dependency{dep1.ID, quantity.Quantity{500, "g"}, 0}, // 50
		Dependency{dep2.ID, quantity.Quantity{1, "kg"}, 0},  // 50
		Dependency{dep3.ID, quantity.Quantity{200, "g"}, 0}, // 25
	}

	repo.InsertMany([]*Composition{dep1, dep2, dep3})

	createReq := compToCreateRequest(comp)
	comp, err := serv.Create(createReq)
	tests.Ok(t, err, "No error")

	tests.Equal(t, comp.Cost, 125.0, "Cost should be 125.0")

	tests.Assert(t, comp.Dependencies[0].Subvalue == 50 && comp.Dependencies[1].Subvalue == 50 && comp.Dependencies[2].Subvalue == 25, "Dependency subvalue wrong")

	t.Run("Add dependency", func(t *testing.T) {
		dep4 := newComposition()
		dep4.Cost = 25
		dep4.Unit = quantity.Quantity{0.5, "u"}
		repo.Insert(dep4)

		q := quantity.Quantity{1, "u"}
		comp.Dependencies = append(comp.Dependencies, Dependency{dep4.ID, q, 0}) // 50

		updateReq := compToUpdateRequest(comp)
		comp, err := serv.Update(comp.ID.Hex(), updateReq)
		tests.Ok(t, err, "No error")

		tests.Equal(t, comp.Cost, 175.0, "Cost should be 175.0")
	})

	t.Run("Remove dependency", func(t *testing.T) {
		// Remove first dependencuy
		comp.Dependencies = comp.Dependencies[1:]

		updateReq := compToUpdateRequest(comp)
		comp, err := serv.Update(comp.ID.Hex(), updateReq)
		tests.Ok(t, err, "No error")
		tests.Equal(t, comp.Cost, 125.0, "Cost should be 125.0")
	})

	t.Run("Change dependency", func(t *testing.T) {
		// Change 1 kg to 8 kg = $400
		comp.Dependencies[0].Quantity = quantity.Quantity{8, "kg"}

		updateReq := compToUpdateRequest(comp)
		comp, err := serv.Update(comp.ID.Hex(), updateReq)
		tests.Ok(t, err, "No error")
		tests.Equal(t, comp.Cost, 475.0, "Cost should be 475.0")
	})
}

func TestDeleteComposition(t *testing.T) {
	repo, eventMgr := newMockRepository(), events.GetMockManager()
	serv := NewService(repo, eventMgr)

	comp, dep := newComposition(), newComposition()
	dep.Cost = 10
	dep.Unit = quantity.Quantity{1, "u"}
	repo.Insert(dep)
	comp.Dependencies = []Dependency{
		Dependency{
			On:       dep.ID,
			Quantity: quantity.Quantity{2, "u"},
		},
	}
	repo.Insert(comp)

	tests.ErrCode(t, serv.Delete(dep.ID.Hex()), "COMPOSITION_USED_AS_DEPENDENCY", "Used dependency cannot be deleted")
	tests.Equal(t, repo.Count(), 2, "Used dependency cannot be deleted")
	tests.Ok(t, serv.Delete(comp.ID.Hex()), "Not used composition should be deleted")
	tests.Equal(t, repo.Count(), 1, "Not used composition should be deleted")
}

func TestCalculateDependenciesSubvalues(t *testing.T) {
	repo, eventMgr := newMockRepository(), events.GetMockManager()
	serv := NewService(repo, eventMgr)

	comps := makeMockedCompositions()
	repo.InsertMany(comps)
	for _, c := range comps {
		servImpl := serv.(*service)
		deps, err := servImpl.calculateDependenciesSubvalues(c.Dependencies)
		tests.Ok(t, err, "No error")
		c.SetDependencies(deps)
		tests.Ok(t, repo.Update(c), "No error")
	}

	c1 := 200.0
	q1 := 2.0
	c2 := 0.2 * c1 / q1 // 20
	c3 := 0.1 * c1 / q1 // 10
	c4 := 150.0
	c5 := 0.4*c2/0.2 + 0.05*c3/0.5 // 41
	c6 := 0.35 * c4 / 0.1          // 525
	c7 := 2*c5/1 + 1.5*c6/2        // 475.75

	checkCompCost(t, comps, 0, c1)
	checkCompCost(t, comps, 1, c2)
	checkCompCost(t, comps, 2, c3)
	checkCompCost(t, comps, 3, c4)
	checkCompCost(t, comps, 4, c5)
	checkCompCost(t, comps, 5, c6)
	checkCompCost(t, comps, 6, c7)
}
