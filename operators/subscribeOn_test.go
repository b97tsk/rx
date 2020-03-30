package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestSubscribeOn(t *testing.T) {
	Subscribe(
		t,
		rx.Merge(
			rx.Just("A", "B").Pipe(operators.SubscribeOn(Step(1))),
			rx.Just("C", "D").Pipe(operators.SubscribeOn(Step(2))),
			rx.Just("E", "F").Pipe(operators.SubscribeOn(Step(3))),
		),
		"A", "B", "C", "D", "E", "F", rx.Complete,
	)
}
