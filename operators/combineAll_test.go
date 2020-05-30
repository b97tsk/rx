package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestCombineAll(t *testing.T) {
	Subscribe(
		t,
		rx.Just(
			rx.Just("A", "B").Pipe(AddLatencyToValues(3, 5)),
			rx.Just("C", "D").Pipe(AddLatencyToValues(2, 4)),
			rx.Just("E", "F").Pipe(AddLatencyToValues(1, 3)),
		).Pipe(operators.CombineAll(), ToString()),
		"[A C E]", "[A C F]", "[A D F]", "[B D F]", rx.Completed,
	)
}
