package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Distinct(t *testing.T) {
	subscribe(
		t,
		Just("A", "B", "B", "A", "C", "C", "A").Pipe(operators.Distinct()),
		"A", "B", "C", Complete,
	)
}
