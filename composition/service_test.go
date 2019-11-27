package composition

import (
	"math"
	"testing"

	"github.com/aboglioli/big-brother/errors"
	"github.com/aboglioli/big-brother/events"
	"github.com/aboglioli/big-brother/quantity"
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
			t.Errorf("- dep %s subvalue %.2f", dep.Of.Hex(), dep.Subvalue)
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

	// Errors
	t.Run("Negative cost", func(t *testing.T) {
		comp := newComposition()
		comp.Cost = -1.0
		_, err := serv.Create(compToCreateRequest(comp))

		if !hasErrCode(err, "NEGATIVE_COST") {
			t.Error("Cost can't be negative")
		}
	})

	t.Run("Invalid units", func(t *testing.T) {
		comp := newComposition()
		comp.Unit.Unit = "asd"

		if _, err := serv.Create(compToCreateRequest(comp)); !hasErrCode(err, "INVALID_UNIT") {
			t.Error("Unit should exist")
		}

		comp.Unit.Unit = "u"
		comp.Stock.Unit = "asd"
		if _, err := serv.Create(compToCreateRequest(comp)); !hasErrCode(err, "INVALID_STOCK") {
			t.Error("Stock unit should exist")
		}

		comp.Unit.Unit = "kg"
		comp.Stock.Unit = "l"
		if _, err := serv.Create(compToCreateRequest(comp)); !hasErrCode(err, "INCOMPATIBLE_STOCK_AND_UNIT") {
			t.Error("Stock and unit should be compatible")
		}
	})

	t.Run("Non-existing dependency", func(t *testing.T) {
		comp := newComposition()
		comp.Dependencies = []Dependency{
			Dependency{
				Of: primitive.NewObjectID(),
				Quantity: quantity.Quantity{
					Quantity: 5,
					Unit:     "u",
				},
			},
		}
		_, err := serv.Create(compToCreateRequest(comp))

		if !hasErrCode(err, "DEPENDENCY_DOES_NOT_EXIST") {
			t.Error("Check dependency existence")
		}
	})

	t.Run("Incompatible dependency quantity", func(t *testing.T) {
		repo.Clean()
		dep, comp := newComposition(), newComposition()
		dep.Unit = quantity.Quantity{2, "kg"}
		repo.Insert(dep)
		comp.Dependencies = []Dependency{
			Dependency{
				Of: dep.ID,
				Quantity: quantity.Quantity{
					Quantity: 5,
					Unit:     "l",
				},
			},
		}

		if _, err := serv.Create(compToCreateRequest(comp)); !hasErrCode(err, "INCOMPATIBLE_DEPENDENCY_QUANTITY") {
			t.Error("Dependency cannot be created with incompatible dependency quantity")
		}
	})

	t.Run("Invalid dependency quantity", func(t *testing.T) {
		repo.Clean()
		dep, comp := newComposition(), newComposition()
		dep.Unit = quantity.Quantity{2, "kg"}
		repo.Insert(dep)
		comp.Dependencies = []Dependency{
			Dependency{
				Of: dep.ID,
				Quantity: quantity.Quantity{
					Quantity: 5,
					Unit:     "kk",
				},
			},
		}

		if _, err := serv.Create(compToCreateRequest(comp)); !hasErrCode(err, "INVALID_DEPENDENCY_QUANTITY") {
			t.Error("Dependency cannot be created with invalid dependency quantity")
		}

		comp.Dependencies[0].Quantity = quantity.Quantity{-5, "kg"}
		if _, err := serv.Create(compToCreateRequest(comp)); !hasErrCode(err, "INVALID_DEPENDENCY_QUANTITY") {
			t.Error("Dependency cannot be created with invalid dependency quantity")
		}
	})

	// Create
	t.Run("Default values with valid units and raise event 'CompositionCreated'", func(t *testing.T) {
		repo.Clean()
		eventMgr.Clean()
		comp := newComposition()
		_, err := serv.Create(compToCreateRequest(comp))

		if err != nil || repo.Count() != 1 {
			t.Error("Composition should be created")
		}

		if eventMgr.Count() != 1 {
			t.Errorf("Creation should emit an event: %d\n", eventMgr.Count())
		}

		msgs := eventMgr.Messages()
		msg := msgs[0]
		if msg.Type() != "CompositionCreated" {
			t.Errorf("Wrong event: %s\n", msg.Type())
		}
	})

	t.Run("Assign stock automatically from unit", func(t *testing.T) {
		repo.Clean()
		comp := newComposition()
		comp.Unit = quantity.Quantity{2, "kg"}
		createReq := compToCreateRequest(comp)
		createReq.Stock = nil
		c, err := serv.Create(createReq)
		if err != nil {
			t.Error(err)
		}

		if !c.Stock.Equals(quantity.Quantity{0, c.Unit.Unit}) {
			t.Error("Stock should be auto-assigned")
		}
	})

	t.Run("Valid dependency", func(t *testing.T) {
		repo.Clean()
		dep, comp := newComposition(), newComposition()
		repo.Insert(dep)
		comp.Dependencies = []Dependency{
			Dependency{
				Of: dep.ID,
				Quantity: quantity.Quantity{
					Quantity: 5,
					Unit:     "u",
				},
			},
		}

		_, err := serv.Create(compToCreateRequest(comp))

		if err != nil || repo.Count() != 2 {
			t.Error(err)
			t.Error("Composition with single dependency should be created")
		}
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
				Of: dep.ID,
				Quantity: quantity.Quantity{
					Quantity: 750,
					Unit:     "g",
				},
			},
		}
		repo.Insert(dep)

		c, err := serv.Create(compToCreateRequest(comp))

		if err != nil {
			t.Error(err)
			return
		}

		if c.Cost != 37.5 {
			t.Error("Cost wrong calculated")
		}

		if eventMgr.Count() != 1 {
			t.Errorf("Should raise an event: %d\n", eventMgr.Count())
		}

		msgs := eventMgr.Messages()
		msg := msgs[0]
		if msg.Type() != "CompositionCreated" {
			t.Errorf("Wrong event: %s\n", msg.Type())
		}

		var evt CompositionChangedEvent
		if err := msg.Decode(&evt); err != nil {
			t.Error(err)
		}
		if evt.Composition.Cost != c.Cost || evt.Composition.ID.Hex() != c.ID.Hex() {
			t.Error("Composition from event is not the expected one")
		}
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
			if err != nil {
				t.Error(err)
				continue
			}
			c.SetDependencies(deps)
			if err := repo.Update(c); err != nil {
				t.Error(err)
				continue
			}
		}

		c := comps[0]
		c.Cost = 300
		c.Unit = quantity.Quantity{
			Quantity: 2500,
			Unit:     "g",
		}

		c, err := serv.Update(c.ID.Hex(), compToUpdateRequest(c))
		if err != nil {
			t.Error(err)
		}
		updatedUses, err := serv.UpdateUses(c)
		if err != nil {
			t.Error(err)
		}

		if len(updatedUses) != 4 {
			t.Errorf("Uses weren't updated: %d updated instead of %d\n", len(updatedUses), 4)
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
		if eventMgr.Count() != 2 {
			t.Errorf("Should raise events: %d\n", eventMgr.Count())
		}

		msgs := eventMgr.Messages()
		if msgs[0].Type() != "CompositionUpdatedManually" {
			t.Errorf("Wrong first event: %s\n", msgs[0].Type())
		}
		var compUpdatedManuallyEvent CompositionChangedEvent
		if err := msgs[0].Decode(&compUpdatedManuallyEvent); err != nil {
			t.Error(err)
		}
		if compUpdatedManuallyEvent.Composition.ID.Hex() != c.ID.Hex() {
			t.Errorf("Different composition in raised event\n")
		}

		if msgs[1].Type() != "CompositionsUpdatedAutomatically" {
			t.Errorf("Wrong first event: %s\n", msgs[0].Type())
		}
		var compsUpdatedAutomaticallyEvent CompositionUpdatedAutomaticallyEvent
		if err := msgs[1].Decode(&compsUpdatedAutomaticallyEvent); err != nil {
			t.Error(err)
		}
		if len(compsUpdatedAutomaticallyEvent.Compositions) != 4 {
			t.Errorf("Compositions updated automatically: %d\n", len(compsUpdatedAutomaticallyEvent.Compositions))
		}
	})

	t.Run("Creating and updating", func(t *testing.T) {
		repo.Clean()
		comp := newComposition()
		comp.Cost = 30
		comp.Unit = quantity.Quantity{1, "u"}
		comp.Stock = quantity.Quantity{1, "u"}

		createdComp, err := serv.Create(compToCreateRequest(comp))
		if err != nil {
			t.Error(err)
		}

		updatedComp, err := serv.Update(createdComp.ID.Hex(), compToUpdateRequest(createdComp))
		if err != nil {
			t.Error(err)
		}

		if createdComp.ID.Hex() != updatedComp.ID.Hex() {
			t.Error("Composition ID has changed")
		}
	})

	t.Run("Invalid units", func(t *testing.T) {
		repo.Clean()
		comp := newComposition()
		comp.Cost = 30
		comp.Unit = quantity.Quantity{1, "u"}
		comp.Stock = quantity.Quantity{1, "u"}

		createdComp, err := serv.Create(compToCreateRequest(comp))
		if err != nil {
			t.Error(err)
		}

		createdComp.Unit = quantity.Quantity{1, "asd"}
		if _, err := serv.Update(createdComp.ID.Hex(), compToUpdateRequest(createdComp)); !hasErrCode(err, "INVALID_UNIT") {
			t.Error("Should return error due to invalid unit")
		}

		createdComp.Unit = quantity.Quantity{1, "u"}
		createdComp.Stock = quantity.Quantity{1, "asd"}
		if _, err := serv.Update(createdComp.ID.Hex(), compToUpdateRequest(createdComp)); !hasErrCode(err, "INVALID_STOCK") {
			t.Error("Should return error due to invalid unit in stock")
		}

		createdComp.Unit = quantity.Quantity{1, "kg"}
		createdComp.Stock = quantity.Quantity{1, "l"}
		if _, err := serv.Update(createdComp.ID.Hex(), compToUpdateRequest(createdComp)); !hasErrCode(err, "INCOMPATIBLE_STOCK_AND_UNIT") {
			t.Error("Stock and unit should be compatible")
		}
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
		if err != nil {
			t.Error(err)
		}

		if !c.Stock.Equals(quantity.Quantity{25, "l"}) {
			t.Error("Empty stock should be ignored")
		}

		updateReq.Stock = &quantity.Quantity{4000, "ml"}
		c, err = serv.Update(comp.ID.Hex(), updateReq)
		if err != nil {
			t.Error(err)
		}

		if !c.Stock.Equals(quantity.Quantity{4000, "ml"}) {
			t.Error("Non-empty stock should be assigned")
		}
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
	if err != nil {
		t.Error(err)
	}

	if comp.Cost != 125 {
		t.Errorf("Cost %f should be %f\n", comp.Cost, 125.0)
	}

	if comp.Dependencies[0].Subvalue != 50 || comp.Dependencies[1].Subvalue != 50 || comp.Dependencies[2].Subvalue != 25 {
		t.Errorf("Dependency subvalue wrong")
	}

	t.Run("Add dependency", func(t *testing.T) {
		dep4 := newComposition()
		dep4.Cost = 25
		dep4.Unit = quantity.Quantity{0.5, "u"}
		repo.Insert(dep4)

		q := quantity.Quantity{1, "u"}
		comp.Dependencies = append(comp.Dependencies, Dependency{dep4.ID, q, 0}) // 50

		updateReq := compToUpdateRequest(comp)
		comp, err := serv.Update(comp.ID.Hex(), updateReq)
		if err != nil {
			t.Error(err)
		}

		if comp.Cost != 175 {
			t.Errorf("Cost %f should be 175\n", comp.Cost)
		}
	})

	t.Run("Remove dependency", func(t *testing.T) {
		// Remove first dependencuy
		comp.Dependencies = comp.Dependencies[1:]

		updateReq := compToUpdateRequest(comp)
		comp, err := serv.Update(comp.ID.Hex(), updateReq)
		if err != nil {
			t.Error(err)
		}

		if comp.Cost != 125 {
			t.Errorf("Cost %f should be 125\n", comp.Cost)
		}
	})

	t.Run("Change dependency", func(t *testing.T) {
		// Change 1 kg to 8 kg = $400
		comp.Dependencies[0].Quantity = quantity.Quantity{8, "kg"}

		updateReq := compToUpdateRequest(comp)
		comp, err := serv.Update(comp.ID.Hex(), updateReq)
		if err != nil {
			t.Error(err)
		}

		if comp.Cost != 475 {
			t.Errorf("Cost %f should be 125\n", comp.Cost)
		}
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
			Of:       dep.ID,
			Quantity: quantity.Quantity{2, "u"},
		},
	}
	repo.Insert(comp)

	if err := serv.Delete(dep.ID.Hex()); !hasErrCode(err, "COMPOSITION_USED_AS_DEPENDENCY") || repo.Count() != 2 {
		t.Error("Used dependency cannot be deleted")
	}

	if err := serv.Delete(comp.ID.Hex()); err != nil || repo.Count() != 1 {
		t.Error("Not used composition should be deleted")
	}
}

func TestCalculateDependenciesSubvalues(t *testing.T) {
	repo, eventMgr := newMockRepository(), events.GetMockManager()
	serv := NewService(repo, eventMgr)

	comps := makeMockedCompositions()
	repo.InsertMany(comps)
	for _, c := range comps {
		servImpl := serv.(*service)
		deps, err := servImpl.calculateDependenciesSubvalues(c.Dependencies)
		if err != nil {
			t.Error(err)
			continue
		}
		c.SetDependencies(deps)
		repo.Update(c)
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
