package quantity

import "testing"

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

		if qTotal.Quantity != 2.5 {
			t.Errorf("Result value: %f", qTotal.Quantity)
		}

		if qTotal.Unit != "kg" {
			t.Errorf("Result unit: %s", qTotal.Unit)
		}
	})

	t.Run("Incompatible units", func(t *testing.T) {
		_, err := q1.Add(q3)

		if err == nil {
			t.Errorf("Should return error")
		}
	})
}

func TestSubtract(t *testing.T) {
	t.Run("Successful", func(t *testing.T) {
		qTotal, _ := q1.Subtract(q2)

		if qTotal.Quantity != 1.5 {
			t.Errorf("Result value: %f", qTotal.Quantity)
		}

		if qTotal.Unit != "kg" {
			t.Errorf("Result unit: %s", qTotal.Unit)
		}
	})

	t.Run("Incompatible units", func(t *testing.T) {
		_, err := q1.Subtract(q3)

		if err == nil {
			t.Errorf("Should return error")
		}
	})
}
