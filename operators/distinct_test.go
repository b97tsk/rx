package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestDistinct(t *testing.T) {
	Subscribe(
		t,
		rx.Just("A", "B", "B", "A", "C", "C", "A").Pipe(operators.Distinct()),
		"A", "B", "C", rx.Complete,
	)
}
