package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestConcatSync(t *testing.T) {
	Subscribe(
		t,
		rx.Just(
			rx.Just("A", "B").Pipe(AddLatencyToValues(3, 5)),
			rx.Just("C", "D").Pipe(AddLatencyToValues(2, 4)),
			rx.Just("E", "F").Pipe(AddLatencyToValues(1, 3)),
		).Pipe(operators.ConcatSyncAll()),
		"A", "B", "C", "D", "E", "F", Completed,
	)
	Subscribe(
		t,
		rx.Timer(Step(1)).Pipe(
			operators.ConcatSyncMapTo(rx.Just("A")),
		),
		"A", Completed,
	)
	Subscribe(
		t,
		rx.Just(rx.Throw(ErrTest)).Pipe(
			operators.ConcatSyncAll(),
		),
		ErrTest,
	)
	Subscribe(
		t,
		rx.Throw(ErrTest).Pipe(
			operators.ConcatSyncAll(),
		),
		ErrTest,
	)
}
