package composition

import (
	"testing"

	"github.com/aboglioli/big-brother/quantity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCalculateCostFromSubvalue(t *testing.T) {
	c := newComposition()
	c.Dependencies = append(
		c.Dependencies,
		&Dependency{
			Subvalue: 100,
		},
		&Dependency{
			Subvalue: 250,
		},
	)
	c.CalculateCost()
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
}

func TestAddAndRemoveCompositionDependencies(t *testing.T) {
	c := newComposition()
	randID := primitive.NewObjectID()

	t.Run("New dependency", func(t *testing.T) {
		c.UpsertDependency(&Dependency{
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
		c.UpsertDependency(&Dependency{
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
		c.UpsertDependency(&Dependency{
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
		err := c.RemoveDependency(randID.String())

		if err != nil || len(c.Dependencies) != 1 {
			t.Error("Dependency should be removed")
		}
	})
	t.Run("Remove non-existing dependency", func(t *testing.T) {
		err := c.RemoveDependency(primitive.NewObjectID().String())

		if err == nil {
			t.Error("Shouldn't be removed")
		}
	})
}

func TestCompareDependencies(t *testing.T) {
	c1, c2 := makeMockedCompositions()[1], makeMockedCompositions()[1]

	c2.Dependencies[0].Of = c1.Dependencies[0].Of

	// With itself
	left, common, right := c1.CompareDependencies(c1)
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
	left, common, right = c1.CompareDependencies(c2)
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

	left, common, right = c1.CompareDependencies(c2)
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
	c1.UpsertDependency(&dep)
	c2.UpsertDependency(&dep)

	left, common, right = c1.CompareDependencies(c2)
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
