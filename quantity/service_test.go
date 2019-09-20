package quantity

import "testing"

var q1, q2, q3 *Quantity
var s Service

func init() {
	q1 = &Quantity{
		Unit:  "kg",
		Value: 2,
	}

	q2 = &Quantity{
		Unit:  "g",
		Value: 500,
	}

	q3 = &Quantity{
		Unit:  "km",
		Value: 3,
	}

	s = NewService()
}

func TestAddSuccessful(t *testing.T) {
	t.Run("Successful", func(t *testing.T) {
		qTotal, _ := s.Add(q1, q2)

		if qTotal.Value != 2.5 {
			t.Errorf("Result value: %f", qTotal.Value)
		}

		if qTotal.Unit != "kg" {
			t.Errorf("Result unit: %s", qTotal.Unit)
		}
	})

	t.Run("Incompatible units", func(t *testing.T) {
		_, err := s.Add(q1, q3)

		if err == nil {
			t.Errorf("Should return error")
		}
	})
}

func TestSubstract(t *testing.T) {
	t.Run("Successful", func(t *testing.T) {
		qTotal, _ := s.Substract(q1, q2)

		if qTotal.Value != 1.5 {
			t.Errorf("Result value: %f", qTotal.Value)
		}

		if qTotal.Unit != "kg" {
			t.Errorf("Result unit: %s", qTotal.Unit)
		}
	})

	t.Run("Incompatible units", func(t *testing.T) {
		_, err := s.Add(q1, q3)

		if err == nil {
			t.Errorf("Should return error")
		}
	})
}
