package quantity

import (
	"testing"

	"github.com/aboglioli/big-brother/pkg/tests/assert"
)

var q1, q2, q3 Quantity

func init() {
	q1 = Quantity{
		Unit:     "kg",
		Quantity: 2,
	}

	q2 = Quantity{
		Unit:     "g",
		Quantity: 500,
	}

	q3 = Quantity{
		Unit:     "km",
		Quantity: 3,
	}
}

func TestAddSuccessful(t *testing.T) {
	t.Run("Successful", func(t *testing.T) {
		qTotal, _ := q1.Add(q2)
		assert.Equal(t, qTotal.Quantity, 2.5)
		assert.Equal(t, qTotal.Unit, "kg")
	})

	t.Run("Incompatible units", func(t *testing.T) {
		_, err := q1.Add(q3)
		assert.Err(t, err)
	})
}

func TestSubtract(t *testing.T) {
	t.Run("Successful", func(t *testing.T) {
		qTotal, _ := q1.Subtract(q2)
		assert.Equal(t, qTotal.Quantity, 1.5)
		assert.Equal(t, qTotal.Unit, "kg")
	})

	t.Run("Incompatible units", func(t *testing.T) {
		_, err := q1.Subtract(q3)
		assert.Err(t, err)
	})
}
