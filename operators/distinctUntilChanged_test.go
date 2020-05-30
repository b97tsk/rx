package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestDistinctUntilChanged(t *testing.T) {
	Subscribe(
		t,
		rx.Just("A", "B", "B", "A", "C", "C", "A").Pipe(operators.DistinctUntilChanged()),
		"A", "B", "A", "C", "A", rx.Completed,
	)
}
