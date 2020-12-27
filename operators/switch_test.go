package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestSwitch(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Just(
			rx.Just("A", "B", "C", "D").Pipe(AddLatencyToValues(0, 2)),
			rx.Just("E", "F", "G", "H").Pipe(AddLatencyToValues(0, 3)),
			rx.Just("I", "J", "K", "L").Pipe(AddLatencyToValues(0, 2)),
		).Pipe(
			AddLatencyToValues(0, 5),
			operators.SwitchAll(),
		),
		"A", "B", "C", "E", "F", "I", "J", "K", "L", Completed,
	).Case(
		rx.Just(
			rx.Just("A", "B", "C", "D").Pipe(AddLatencyToValues(0, 2)),
			rx.Just("E", "F", "G", "H").Pipe(AddLatencyToValues(0, 3)),
			rx.Just("I", "J", "K", "L").Pipe(AddLatencyToValues(0, 2)),
			rx.Empty(),
		).Pipe(
			AddLatencyToValues(0, 5),
			operators.SwitchAll(),
		),
		"A", "B", "C", "E", "F", "I", "J", "K", Completed,
	).Case(
		rx.Just(
			rx.Just("A", "B", "C", "D").Pipe(AddLatencyToValues(0, 2)),
			rx.Just("E", "F", "G", "H").Pipe(AddLatencyToValues(0, 3)),
			rx.Just("I", "J", "K", "L").Pipe(AddLatencyToValues(0, 2)),
			rx.Throw(ErrTest),
		).Pipe(
			AddLatencyToValues(0, 5),
			operators.SwitchAll(),
		),
		"A", "B", "C", "E", "F", "I", "J", "K", ErrTest,
	).Case(
		rx.Just("A").Pipe(
			operators.SwitchMapTo(rx.Empty()),
		),
		Completed,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.SwitchAll(),
		),
		ErrTest,
	).TestAll()
}
