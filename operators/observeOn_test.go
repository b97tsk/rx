package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestObserveOn(t *testing.T) {
	Subscribe(
		t,
		rx.Merge(
			rx.Just("A", "B").Pipe(operators.ObserveOn(Step(1))),
			rx.Just("C", "D").Pipe(operators.ObserveOn(Step(2))),
			rx.Just("E", "F").Pipe(operators.ObserveOn(Step(3))),
		),
		"A", "B", "C", "D", "E", "F", rx.Complete,
	)
}
