package composition

import (
	"testing"

	"github.com/aboglioli/big-brother/pkg/quantity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCalculateCostFromSubvalue(t *testing.T) {
	c := newComposition()
	c.Dependencies = append(
		c.Dependencies,
		Dependency{
			Subvalue: 100,
		},
		Dependency{
			Subvalue: 250,
		},
	)
	c.calculateCostFromDependencies()
	if c.Cost != 350 {
		t.Error("Cost should be 350")
	}
}

func TestCalculateCostByQuantity(t *testing.T) {
	comp := NewComposition()
	comp.Cost = 50
	comp.Unit = quantity.Quantity{2, "kg"}

	if comp.CostFromQuantity(quantity.Quantity{1000, "g"}) != 25 {
		t.Error("Cost should be 25")
	}

	if comp.CostFromQuantity(quantity.Quantity{500, "g"}) != 12.5 {
		t.Error("Cost should be 12.5")
	}

	if comp.CostFromQuantity(quantity.Quantity{3, "kg"}) != 3*50/2 {
		t.Error("Cost should be 12.5")
	}

	comp.Unit = quantity.Quantity{0, "kg"}

	if comp.CostFromQuantity(quantity.Quantity{1, "kg"}) != 0 {
		t.Error("Division by zero")
	}
}

func TestAddAndRemoveCompositionDependencies(t *testing.T) {
	c := newComposition()
	randID := primitive.NewObjectID()

	t.Run("New dependency", func(t *testing.T) {
		c.UpsertDependency(Dependency{
			Of: randID,
			Quantity: quantity.Quantity{
				Unit:     "u",
				Quantity: 1,
			},
		})

		if len(c.Dependencies) != 1 {
			t.Error("Dependency should have been added")
		}
	})
	t.Run("Add same dependency", func(t *testing.T) {
		c.UpsertDependency(Dependency{
			Of: randID,
			Quantity: quantity.Quantity{
				Unit:     "u",
				Quantity: 2,
			},
		})

		if len(c.Dependencies) != 1 || c.Dependencies[0].Quantity.Quantity != 2 {
			t.Error("Upsert same dependency")
		}
	})
	t.Run("Add new dependency", func(t *testing.T) {
		c.UpsertDependency(Dependency{
			Of: primitive.NewObjectID(),
			Quantity: quantity.Quantity{
				Unit:     "u",
				Quantity: 1.5,
			},
		})

		if len(c.Dependencies) != 2 || c.Dependencies[1].Quantity.Quantity != 1.5 {
			t.Error("Upsert different dependency")
		}
	})
	t.Run("Remove existing dependency", func(t *testing.T) {
		err := c.RemoveDependency(randID.Hex())

		if err != nil || len(c.Dependencies) != 1 {
			t.Error("Dependency should be removed")
		}
	})
	t.Run("Remove non-existing dependency", func(t *testing.T) {
		err := c.RemoveDependency(primitive.NewObjectID().String())

		if err == nil {
			t.Error("Should throw an error")
		}
	})
	t.Run("Add dependencies and calculate cost from subvalues", func(t *testing.T) {
		c.Dependencies = []Dependency{}
		c.UpsertDependency(Dependency{
			Of: primitive.NewObjectID(),
			Quantity: quantity.Quantity{
				Unit:     "u",
				Quantity: 1.5,
			},
			Subvalue: 20.5,
		})
		id := primitive.NewObjectID()
		c.UpsertDependency(Dependency{
			Of: id,
			Quantity: quantity.Quantity{
				Unit:     "u",
				Quantity: 2.5,
			},
			Subvalue: 30,
		})
		c.UpsertDependency(Dependency{
			Of: primitive.NewObjectID(),
			Quantity: quantity.Quantity{
				Unit:     "u",
				Quantity: 2.5,
			},
			Subvalue: 10.5,
		})

		if len(c.Dependencies) != 3 || c.Cost != 61 {
			t.Error("Cost must be calculate after upserting")
		}

		c.RemoveDependency(id.Hex())

		if len(c.Dependencies) != 2 || c.Cost != 31 {
			t.Error("Cost must be calculate after removing")
		}
	})

	t.Run("Add new dependency to a non-autoupdated composition", func(t *testing.T) {
		c := newComposition()
		c.Cost = 45
		c.Unit = quantity.Quantity{2, "kg"}
		c.Stock = c.Unit
		c.AutoupdateCost = false

		c.UpsertDependency(Dependency{
			Of:       primitive.NewObjectID(),
			Quantity: quantity.Quantity{1.5, "u"},
			Subvalue: 30,
		})

		if c.Cost != 45 {
			t.Error("Cost should not be updated automatically")
		}
	})
}

func TestCompareDependencies(t *testing.T) {
	c1, c2 := makeMockedCompositions()[1], makeMockedCompositions()[1]

	c2.Dependencies[0].Of = c1.Dependencies[0].Of

	// With itself
	left, common, right := c1.CompareDependencies(c1.Dependencies)
	if len(left) != 0 {
		t.Error("Left should be empty")
	}
	if len(common) != 1 {
		t.Error("Common should have 1 element")
	}
	if len(right) != 0 {
		t.Error("Right should be empty")
	}

	// Copy
	left, common, right = c1.CompareDependencies(c2.Dependencies)
	if len(left) != 0 {
		t.Error("Left should be empty")
	}
	if len(common) != 1 {
		t.Error("Common should have 1 element")
	}
	if len(right) != 0 {
		t.Error("Right should be empty")
	}

	// After changing
	c2.Dependencies[0].Quantity = quantity.Quantity{250.0, "g"}

	left, common, right = c1.CompareDependencies(c2.Dependencies)
	if len(left) != 1 {
		t.Error("Left")
	}
	if len(common) != 0 {
		t.Error("There aren't common dependencies")
	}
	if len(right) != 1 {
		t.Error("Right")
	}

	// Add a common dependency
	dep := Dependency{
		Of:       primitive.NewObjectID(),
		Quantity: quantity.Quantity{2, "l"},
	}
	c1.UpsertDependency(dep)
	c2.UpsertDependency(dep)

	left, common, right = c1.CompareDependencies(c2.Dependencies)
	if len(left) != 1 {
		t.Error("Left")
	}
	if len(common) != 1 {
		t.Error("There is a common dependency")
	}
	if len(right) != 1 {
		t.Error("Right")
	}
}
