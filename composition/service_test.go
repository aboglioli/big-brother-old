package composition

import (
	"math"
	"testing"

	"github.com/aboglioli/big-brother/impl/events"
	"github.com/aboglioli/big-brother/pkg/quantity"
	"github.com/aboglioli/big-brother/pkg/tests/assert"
	"github.com/aboglioli/big-brother/pkg/tests/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
	return &UpdateRequest{
		Name:           &c.Name,
		Cost:           &c.Cost,
		Unit:           &c.Unit,
		Stock:          &c.Stock,
		Dependencies:   c.Dependencies,
		AutoupdateCost: &c.AutoupdateCost,
	}
}

func TestGetByID(t *testing.T) {
	repo, eventMgr := newMockRepository(), events.InMemory()
	serv := NewService(repo, eventMgr)

	// Errors
	t.Run("Not existing", func(t *testing.T) {
		_, err := serv.GetByID("123")
		assert.ErrCode(t, err, "COMPOSITION_NOT_FOUND")
	})

	t.Run("Disabled", func(t *testing.T) {
		comp := newComposition()
		comp.Enabled = false
		comp.Validated = true
		repo.Insert(comp)
		_, err := serv.GetByID(comp.ID.Hex())
		assert.ErrCode(t, err, "COMPOSITION_NOT_FOUND")
	})

	t.Run("Not validated", func(t *testing.T) {
		comp := newComposition()
		comp.Enabled = true
		comp.Validated = false
		repo.Insert(comp)
		_, err := serv.GetByID(comp.ID.Hex())
		assert.ErrCode(t, err, "COMPOSITION_NOT_VALIDATED")
	})

	// OK
	t.Run("OK", func(t *testing.T) {
		comp := newComposition()
		comp.Enabled = true
		comp.Validated = true
		repo.Insert(comp)
		saved, err := serv.GetByID(comp.ID.Hex())
		assert.Ok(t, err)
		assert.Equal(t, saved.ID.Hex(), comp.ID.Hex())
	})
}

func TestCreateComposition(t *testing.T) {
	repo, eventMgr := newMockRepository(), events.InMemory()
	serv := NewService(repo, eventMgr)

	// Errors
	t.Run("Invalid ID", func(t *testing.T) {
		comp := newComposition()
		req := compToCreateRequest(comp)
		id := "123"
		req.ID = &id
		_, err := serv.Create(req)
		assert.ErrCode(t, err, "INVALID_ID")
	})
	t.Run("Existing composition", func(t *testing.T) {
		comp := newComposition()
		repo.Insert(comp)
		_, err := serv.Create(compToCreateRequest(comp))
		assert.ErrCode(t, err, "COMPOSITION_ALREADY_EXISTS")
	})
	t.Run("Existing composition", func(t *testing.T) {
		comp := newComposition()
		repo.Insert(comp)
		_, err := serv.Create(compToCreateRequest(comp))
		assert.ErrCode(t, err, "COMPOSITION_ALREADY_EXISTS")
	})
	t.Run("Non-existing dependency", func(t *testing.T) {
		repo.Mock.Reset()
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
		assert.ErrCode(t, err, "DEPENDENCY_DOES_NOT_EXIST", "Check dependency existence")
		repo.Mock.Assert(t,
			mock.Call("FindByID", comp.ID.Hex()),
			mock.Call("FindByID", comp.Dependencies[0].On.Hex()),
		)
	})

	t.Run("Incompatible dependency quantity", func(t *testing.T) {
		repo.Clean()
		dep, comp := newComposition(), newComposition()
		dep.Unit = quantity.Quantity{2, "kg"}
		repo.Insert(dep)
		repo.Mock.Reset()
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
		assert.ErrCode(t, err, "INCOMPATIBLE_DEPENDENCY_QUANTITY", "Dependency cannot be created with incompatible dependency quantity")
		repo.Mock.Assert(t,
			mock.Call("FindByID", comp.ID.Hex()),
			mock.Call("FindByID", comp.Dependencies[0].On.Hex()),
		)
	})

	// OK
	t.Run("Default values with valid units and raise event 'CompositionCreated'", func(t *testing.T) {
		repo.Clean()
		repo.Mock.Reset()
		eventMgr.Clean()
		eventMgr.Mock.Reset()
		comp := newComposition()

		_, err := serv.Create(compToCreateRequest(comp))
		assert.Ok(t, err)
		total, _ := repo.Count()
		assert.Equal(t, total, 1)
		assert.Equal(t, eventMgr.Count(), 1)

		repo.Mock.Assert(t,
			mock.Call("FindByID", comp.ID.Hex()),
			mock.Call("Insert", mock.NotNil),
		)
		savedComp, ok := repo.Calls[1].Args[0].(*Composition)
		assert.Assert(t, ok)
		assert.Equal(t, savedComp.ID.Hex(), comp.ID.Hex())
		assert.Equal(t, savedComp.Enabled, true)
		assert.Equal(t, savedComp.Validated, false)

		eventMgr.Mock.Assert(t,
			mock.Call("Publish", mock.NotNil, mock.NotNil),
		)
		msgs := eventMgr.Messages()
		msg := msgs[0]
		assert.Equal(t, msg.Type(), "CompositionCreated")
	})

	t.Run("Assign stock automatically from unit", func(t *testing.T) {
		repo.Clean()
		comp := newComposition()
		comp.Unit = quantity.Quantity{2, "kg"}
		createReq := compToCreateRequest(comp)
		createReq.Stock = nil
		c, err := serv.Create(createReq)

		assert.Ok(t, err)
		assert.Assert(t, c.Stock.Equals(quantity.Quantity{0, c.Unit.Unit}))
	})

	t.Run("Valid dependency", func(t *testing.T) {
		repo.Clean()
		dep, comp := newComposition(), newComposition()
		repo.Insert(dep)
		repo.Mock.Reset()
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

		repo.Mock.Assert(t,
			mock.Call("FindByID", comp.ID.Hex()),
			mock.Call("FindByID", dep.ID.Hex()),
			mock.Call("Insert", mock.NotNil),
		)
		savedDepID, ok := repo.Calls[1].Args[0].(string)
		assert.Assert(t, ok)
		savedComp, ok := repo.Calls[2].Args[0].(*Composition)
		assert.Assert(t, ok)
		assert.Equal(t, savedDepID, dep.ID.Hex())
		assert.Equal(t, savedComp.ID.Hex(), comp.ID.Hex())

		assert.Ok(t, err)
		total, _ := repo.Count()
		assert.Equal(t, total, 2)
	})

	t.Run("Calculate cost on creating and raise event", func(t *testing.T) {
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
		eventMgr.Mock.Reset()

		c, err := serv.Create(compToCreateRequest(comp))

		assert.Ok(t, err)
		assert.NotNil(t, c)
		assert.Equal(t, c.Cost, 37.5, "Cost not calculated")
		assert.Equal(t, eventMgr.Count(), 1, "Should raise an event")

		eventMgr.Mock.Assert(t,
			mock.Call("Publish", mock.NotNil, mock.NotNil),
		)
		actualEvent, ok := eventMgr.Calls[0].Args[0].(*CompositionChangedEvent)
		assert.Assert(t, ok)
		assert.Equal(t, actualEvent.Type, "CompositionCreated")
		assert.NotNil(t, actualEvent.Composition)
		assert.Equal(t, actualEvent.Composition.ID.Hex(), comp.ID.Hex())

		msgs := eventMgr.Messages()
		msg := msgs[0]
		assert.Assert(t, msg.Type() == "CompositionCreated" && msg.Key == "composition.created", "Wrong event")

		var evt CompositionChangedEvent
		assert.Ok(t, msg.Decode(&evt), "Decode")
		assert.Assert(t, evt.Composition.Cost == c.Cost && evt.Composition.ID.Hex() == c.ID.Hex(), "Composition from event is not the expected one")
	})
}

func TestUpdateComposition(t *testing.T) {
	repo, eventMgr := newMockRepository(), events.InMemory()
	serv := NewService(repo, eventMgr)

	// Errors
	t.Run("Wrong ID", func(t *testing.T) {
		comp := newComposition()
		repo.Insert(comp)
		repo.Mock.Reset()
		req := compToUpdateRequest(comp)

		_, err := serv.Update("123", req)
		assert.ErrCode(t, err, "COMPOSITION_NOT_FOUND")

		_, err = serv.Update(comp.ID.Hex()+"---", req)
		assert.ErrCode(t, err, "COMPOSITION_NOT_FOUND")

		updatedComp, err := serv.Update(comp.ID.Hex(), req)
		assert.Ok(t, err)
		assert.Equal(t, updatedComp.ID.Hex(), comp.ID.Hex())
	})

	t.Run("Composition disabled and not validated", func(t *testing.T) {
		repo.Clean()
		eventMgr.Clean()
		comp := newComposition()

		comp.Enabled = false
		comp.Validated = true
		repo.Insert(comp)
		_, err := serv.Update(comp.ID.Hex(), compToUpdateRequest(comp))
		assert.ErrCode(t, err, "COMPOSITION_NOT_FOUND")

		comp.Enabled = true
		comp.Validated = false
		repo.Update(comp)
		_, err = serv.Update(comp.ID.Hex(), compToUpdateRequest(comp))
		assert.ErrCode(t, err, "COMPOSITION_NOT_VALIDATED")
	})

	t.Run("Invalid units", func(t *testing.T) {
		repo.Clean()
		comp := newComposition()
		comp.Cost = 30
		comp.Unit = quantity.Quantity{1, "u"}
		comp.Stock = quantity.Quantity{1, "u"}

		createdComp, err := serv.Create(compToCreateRequest(comp))
		assert.Ok(t, err)
		err = serv.Validate(createdComp.ID.Hex())
		assert.Ok(t, err)

		createdComp.Unit = quantity.Quantity{1, "asd"}
		_, err = serv.Update(createdComp.ID.Hex(), compToUpdateRequest(createdComp))
		assert.ErrValidation(t, err, "unit", "INVALID")

		createdComp.Unit = quantity.Quantity{1, "u"}
		createdComp.Stock = quantity.Quantity{1, "asd"}
		_, err = serv.Update(createdComp.ID.Hex(), compToUpdateRequest(createdComp))
		assert.ErrValidation(t, err, "stock", "INVALID")

		createdComp.Unit = quantity.Quantity{1, "kg"}
		createdComp.Stock = quantity.Quantity{1, "l"}
		_, err = serv.Update(createdComp.ID.Hex(), compToUpdateRequest(createdComp))
		assert.ErrValidation(t, err, "stock", "INCOMPATIBLE_STOCK_AND_UNIT")
	})

	t.Run("Change unit after creating", func(t *testing.T) {
		repo.Clean()
		comp := newComposition()
		comp.Unit = quantity.Quantity{1, "kg"}
		comp.Stock = comp.Unit
		repo.Insert(comp)

		comp.Stock.Unit = "l"
		comp.Unit.Unit = "l"
		_, err := serv.Update(comp.ID.Hex(), compToUpdateRequest(comp))
		assert.ErrCode(t, err, "CANNOT_CHANGE_UNIT_TYPE")

		comp.Unit.Unit = "g"
		comp.Stock.Unit = "cg"
		_, err = serv.Update(comp.ID.Hex(), compToUpdateRequest(comp))
		assert.Ok(t, err)
	})

	// OK
	t.Run("Empty stock ignored on updating", func(t *testing.T) {
		repo.Clean()
		comp := newComposition()
		comp.Unit = quantity.Quantity{5, "l"}
		comp.Stock = quantity.Quantity{25, "l"}
		repo.Insert(comp)

		updateReq := compToUpdateRequest(comp)
		updateReq.Stock = nil

		c, err := serv.Update(comp.ID.Hex(), updateReq)
		assert.Ok(t, err)

		assert.Assert(t, c.Stock.Equals(quantity.Quantity{25, "l"}), "Empty stock should be ignored")

		updateReq.Stock = &quantity.Quantity{4000, "ml"}
		c, err = serv.Update(comp.ID.Hex(), updateReq)
		assert.Ok(t, err)

		assert.Assert(t, c.Stock.Equals(quantity.Quantity{4000, "ml"}), "Non-empty stock should be assigned")
	})

	t.Run("Update dependency and raise events", func(t *testing.T) {
		repo.Clean()
		eventMgr.Clean()
		eventMgr.Mock.Reset()

		comps := makeMockedCompositions()
		repo.InsertMany(comps)
		for _, c := range comps {
			servImpl := serv.(*service)
			err := servImpl.validateSchema(c)
			assert.Ok(t, err)
			assert.Ok(t, repo.Update(c))
		}

		c := comps[0]
		c.Cost = 300
		c.Unit = quantity.Quantity{
			Quantity: 2500,
			Unit:     "g",
		}

		c, err := serv.Update(c.ID.Hex(), compToUpdateRequest(c))
		assert.Ok(t, err)
		updatedUses, err := serv.UpdateUses(c)
		assert.Ok(t, err)

		assert.Equal(t, len(updatedUses), 4)

		comps, _ = repo.FindAll()
		assert.Equal(t, len(comps), 7)

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
		eventMgr.Mock.Assert(t,
			mock.Call("Publish", mock.NotNil, mock.NotNil),
			mock.Call("Publish", mock.NotNil, mock.NotNil),
		)
		assert.Equal(t, eventMgr.Count(), 2, "Should raise events")

		msgs := eventMgr.Messages()
		assert.Assert(t, msgs[0].Type() == "CompositionUpdatedManually" && msgs[0].Key == "composition.updated", "Wrong event")

		var compUpdatedManuallyEvent CompositionChangedEvent
		assert.Ok(t, msgs[0].Decode(&compUpdatedManuallyEvent), "Error")
		assert.Equal(t, compUpdatedManuallyEvent.Composition.ID.Hex(), c.ID.Hex(), "Different composition")
		assert.Assert(t, msgs[1].Type() == "CompositionsUpdatedAutomatically" && msgs[0].Key == "composition.updated", "Wrong event")

		var compsUpdatedAutomaticallyEvent CompositionsUpdatedAutomaticallyEvent
		assert.Ok(t, msgs[1].Decode(&compsUpdatedAutomaticallyEvent))
		assert.Equal(t, len(compsUpdatedAutomaticallyEvent.Compositions), 4, "Update automatically")
	})

	t.Run("Creating, validating and updating", func(t *testing.T) {
		repo.Clean()
		comp := newComposition()
		comp.Cost = 30
		comp.Unit = quantity.Quantity{1, "u"}
		comp.Stock = quantity.Quantity{1, "u"}

		createdComp, err := serv.Create(compToCreateRequest(comp))
		assert.Ok(t, err)
		repo.Mock.Reset()
		err = serv.Validate(createdComp.ID.Hex())
		assert.Ok(t, err)

		repo.Mock.Assert(t,
			mock.Call("FindByID", comp.ID.Hex()),
			mock.Call("Update", mock.NotNil),
		)
		validatedComp := repo.Calls[1].Args[0].(*Composition)
		assert.Equal(t, validatedComp.ID.Hex(), comp.ID.Hex())
		assert.Equal(t, validatedComp.Validated, true)

		updatedComp, err := serv.Update(createdComp.ID.Hex(), compToUpdateRequest(createdComp))
		assert.Ok(t, err)
		assert.Equal(t, createdComp.ID.Hex(), updatedComp.ID.Hex(), "ID changed")
	})

}

func TestCreateAndUpdateDependencies(t *testing.T) {
	repo, eventMgr := newMockRepository(), events.InMemory()
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
	assert.Ok(t, err)
	err = serv.Validate(comp.ID.Hex())
	assert.Ok(t, err)

	assert.Equal(t, comp.Cost, 125.0)

	assert.Assert(t, comp.Dependencies[0].Subvalue == 50 && comp.Dependencies[1].Subvalue == 50 && comp.Dependencies[2].Subvalue == 25)

	t.Run("Add dependency", func(t *testing.T) {
		dep4 := newComposition()
		dep4.Cost = 25
		dep4.Unit = quantity.Quantity{0.5, "u"}
		repo.Insert(dep4)

		q := quantity.Quantity{1, "u"}
		comp.Dependencies = append(comp.Dependencies, Dependency{dep4.ID, q, 0}) // 50

		updateReq := compToUpdateRequest(comp)
		comp, err := serv.Update(comp.ID.Hex(), updateReq)
		assert.Ok(t, err)

		assert.Equal(t, comp.Cost, 175.0)
	})

	t.Run("Remove dependency", func(t *testing.T) {
		// Remove first dependencuy
		comp.Dependencies = comp.Dependencies[1:]

		updateReq := compToUpdateRequest(comp)
		comp, err := serv.Update(comp.ID.Hex(), updateReq)
		assert.Ok(t, err)
		assert.Equal(t, comp.Cost, 125.0, "Cost should be 125.0")
	})

	t.Run("Change dependency", func(t *testing.T) {
		// Change 1 kg to 8 kg = $400
		comp.Dependencies[0].Quantity = quantity.Quantity{8, "kg"}

		updateReq := compToUpdateRequest(comp)
		comp, err := serv.Update(comp.ID.Hex(), updateReq)
		assert.Ok(t, err)
		assert.Equal(t, comp.Cost, 475.0, "Cost should be 475.0")
	})
}

func TestDeleteComposition(t *testing.T) {
	repo, eventMgr := newMockRepository(), events.InMemory()
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

	assert.ErrCode(t, serv.Delete(dep.ID.Hex()), "COMPOSITION_USED_AS_DEPENDENCY", "Used dependency cannot be deleted")
	total, enabled := repo.Count()
	assert.Equal(t, total, 2, "Used dependency cannot be deleted")
	assert.Equal(t, total, enabled, "Used dependency cannot be deleted")

	repo.Mock.Reset()
	eventMgr.Mock.Reset()
	assert.Ok(t, serv.Delete(comp.ID.Hex()), "Not used composition should be deleted")
	total, enabled = repo.Count()
	assert.Equal(t, total, 2, "Used dependency cannot be deleted")
	assert.Equal(t, enabled, 1, "Not used composition should be deleted")

	repo.Mock.Assert(t,
		mock.Call("FindByID", comp.ID.Hex()),
		mock.Call("FindUses", comp.ID.Hex()),
		mock.Call("Delete", comp.ID.Hex()),
	)

	eventMgr.Mock.Assert(t,
		mock.Call("Publish", mock.NotNil, mock.NotNil),
	)

	t.Run("Composition disabled and not validated", func(t *testing.T) {
		repo.Clean()
		eventMgr.Clean()

		comp := newComposition()
		comp.Enabled = false
		comp.Validated = true
		repo.Insert(comp)
		err := serv.Delete(comp.ID.Hex())
		assert.ErrCode(t, err, "COMPOSITION_NOT_FOUND")

		comp.Enabled = true
		comp.Validated = false
		repo.Update(comp)
		err = serv.Delete(comp.ID.Hex())
		assert.ErrCode(t, err, "COMPOSITION_NOT_VALIDATED")
	})
}

func TestCalculateDependenciesSubvalues(t *testing.T) {
	repo, eventMgr := newMockRepository(), events.InMemory()
	serv := NewService(repo, eventMgr)

	comps := makeMockedCompositions()
	repo.InsertMany(comps)
	for _, c := range comps {
		servImpl := serv.(*service)
		err := servImpl.validateSchema(c)
		assert.Ok(t, err)
		assert.Ok(t, repo.Update(c))
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
