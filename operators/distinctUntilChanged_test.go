package operators_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_DistinctUntilChanged(t *testing.T) {
	subscribe(
		t,
		Just("A", "B", "B", "A", "C", "C", "A").Pipe(operators.DistinctUntilChanged()),
		"A", "B", "A", "C", "A", Complete,
	)
}
