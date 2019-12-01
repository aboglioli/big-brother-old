package composition

import (
	"testing"

	"github.com/aboglioli/big-brother/pkg/quantity"
	"github.com/aboglioli/big-brother/pkg/tests"
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
	tests.Equal(t, c.Cost, 350.0, "Cost should be 350")
}

func TestCalculateCostByQuantity(t *testing.T) {
	comp := NewComposition()
	comp.Cost = 50
	comp.Unit = quantity.Quantity{2, "kg"}

	tests.Equal(t, comp.CostFromQuantity(quantity.Quantity{1000, "g"}), 25.0, "Cost should be 25")
	tests.Equal(t, comp.CostFromQuantity(quantity.Quantity{500, "g"}), 12.5, "Cost should be 12.5")
	tests.Equal(t, comp.CostFromQuantity(quantity.Quantity{3, "kg"}), 3.0*50/2, "Cost should be 75")

	comp.Unit = quantity.Quantity{0, "kg"}

	tests.Equal(t, comp.CostFromQuantity(quantity.Quantity{1, "kg"}), 0.0, "Division by zero")
}

func TestAddAndRemoveCompositionDependencies(t *testing.T) {
	c := newComposition()
	randID := primitive.NewObjectID()

	t.Run("New dependency", func(t *testing.T) {
		c.UpsertDependency(Dependency{
			On: randID,
			Quantity: quantity.Quantity{
				Unit:     "u",
				Quantity: 1,
			},
		})

		tests.Equal(t, len(c.Dependencies), 1, "Dependency should have been added")
	})
	t.Run("Add same dependency", func(t *testing.T) {
		c.UpsertDependency(Dependency{
			On: randID,
			Quantity: quantity.Quantity{
				Unit:     "u",
				Quantity: 2,
			},
		})

		tests.Assert(t, len(c.Dependencies) == 1 && c.Dependencies[0].Quantity.Quantity == 2, "Upsert same dependency")
	})
	t.Run("Add new dependency", func(t *testing.T) {
		c.UpsertDependency(Dependency{
			On: primitive.NewObjectID(),
			Quantity: quantity.Quantity{
				Unit:     "u",
				Quantity: 1.5,
			},
		})

		tests.Assert(t, len(c.Dependencies) == 2 && c.Dependencies[1].Quantity.Quantity == 1.5, "Upsert different dependency")
	})
	t.Run("Remove existing dependency", func(t *testing.T) {
		err := c.RemoveDependency(randID.Hex())

		tests.Ok(t, err, "Dependency should be removed")
		tests.Equal(t, len(c.Dependencies), 1, "Dependency should be removed")
	})
	t.Run("Remove non-existing dependency", func(t *testing.T) {
		err := c.RemoveDependency(primitive.NewObjectID().String())

		tests.Err(t, err, "Should throw an error")
	})
	t.Run("Add dependencies and calculate cost from subvalues", func(t *testing.T) {
		c.Dependencies = []Dependency{}
		c.UpsertDependency(Dependency{
			On: primitive.NewObjectID(),
			Quantity: quantity.Quantity{
				Unit:     "u",
				Quantity: 1.5,
			},
			Subvalue: 20.5,
		})
		id := primitive.NewObjectID()
		c.UpsertDependency(Dependency{
			On: id,
			Quantity: quantity.Quantity{
				Unit:     "u",
				Quantity: 2.5,
			},
			Subvalue: 30,
		})
		c.UpsertDependency(Dependency{
			On: primitive.NewObjectID(),
			Quantity: quantity.Quantity{
				Unit:     "u",
				Quantity: 2.5,
			},
			Subvalue: 10.5,
		})

		tests.Equal(t, len(c.Dependencies), 3, "Cost should be calculated after upserting")
		tests.Equal(t, c.Cost, 61.0, "Cost should be calculated after upserting")

		c.RemoveDependency(id.Hex())

		tests.Equal(t, len(c.Dependencies), 2, "Cost should be calculated after removing")
		tests.Equal(t, c.Cost, 31.0, "Cost should be calculated after removing")
	})

	t.Run("Add new dependency to a non-autoupdated composition", func(t *testing.T) {
		c := newComposition()
		c.Cost = 45
		c.Unit = quantity.Quantity{2, "kg"}
		c.Stock = c.Unit
		c.AutoupdateCost = false

		c.UpsertDependency(Dependency{
			On:       primitive.NewObjectID(),
			Quantity: quantity.Quantity{1.5, "u"},
			Subvalue: 30,
		})

		tests.Equal(t, c.Cost, 45.0, "Cost should not be updated automatically")
	})
}

func TestCompareDependencies(t *testing.T) {
	c1, c2 := makeMockedCompositions()[1], makeMockedCompositions()[1]

	c2.Dependencies[0].On = c1.Dependencies[0].On

	// With itself
	left, common, right := c1.CompareDependencies(c1.Dependencies)
	tests.Equal(t, len(left), 0, "Left should be empty")
	tests.Equal(t, len(common), 1, "Common should have 1 element")
	tests.Equal(t, len(right), 0, "Right should be empty")

	// Copy
	left, common, right = c1.CompareDependencies(c2.Dependencies)
	tests.Equal(t, len(left), 0, "Left should be empty")
	tests.Equal(t, len(common), 1, "Common should have 1 element")
	tests.Equal(t, len(right), 0, "Right should be empty")

	// After changing
	c2.Dependencies[0].Quantity = quantity.Quantity{250.0, "g"}

	left, common, right = c1.CompareDependencies(c2.Dependencies)
	tests.Equal(t, len(left), 1, "Left")
	tests.Equal(t, len(common), 0, "There aren't common dependencies")
	tests.Equal(t, len(right), 1, "Right")

	// Add a common dependency
	dep := Dependency{
		On:       primitive.NewObjectID(),
		Quantity: quantity.Quantity{2, "l"},
	}
	c1.UpsertDependency(dep)
	c2.UpsertDependency(dep)

	left, common, right = c1.CompareDependencies(c2.Dependencies)
	tests.Equal(t, len(left), 1, "Left")
	tests.Equal(t, len(common), 1, "There is a common dependencies")
	tests.Equal(t, len(right), 1, "Right")
}

func TestValidateSchema(t *testing.T) {
	// Errors
	t.Run("Negative cost", func(t *testing.T) {
		comp := newComposition()
		comp.Cost = -1.0
		tests.ErrCode(t, comp.ValidateSchema(), "NEGATIVE_COST", "Cost should be greater or equal than 0")
	})

	t.Run("Invalid units", func(t *testing.T) {
		comp := newComposition()
		comp.Unit.Unit = "asd"
		tests.ErrCode(t, comp.ValidateSchema(), "INVALID_UNIT", "Unit should exist")

		comp.Unit.Unit = "u"
		comp.Stock.Unit = "asd"
		tests.ErrCode(t, comp.ValidateSchema(), "INVALID_STOCK", "Stock unit should exist")

		comp.Unit.Unit = "kg"
		comp.Stock.Unit = "l"
		tests.ErrCode(t, comp.ValidateSchema(), "INCOMPATIBLE_STOCK_AND_UNIT", "Stock and unit should be compatible")
	})

	t.Run("Invalid dependency quantity", func(t *testing.T) {
		comp := newComposition()
		comp.Dependencies = []Dependency{
			Dependency{
				On: primitive.NewObjectID(),
				Quantity: quantity.Quantity{
					Quantity: 5,
					Unit:     "kk",
				},
			},
		}

		tests.ErrCode(t, comp.ValidateSchema(), "INVALID_DEPENDENCY_QUANTITY", "Dependency has invalid quantity")

		comp.Dependencies[0].Quantity = quantity.Quantity{-5, "kg"}
		tests.ErrCode(t, comp.ValidateSchema(), "INVALID_DEPENDENCY_QUANTITY", "Dependency has invalid quantity")
	})

	// Create
	t.Run("Default values with valid units", func(t *testing.T) {
		comp := newComposition()
		tests.Ok(t, comp.ValidateSchema(), "Should be created")
	})
}
