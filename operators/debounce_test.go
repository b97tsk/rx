package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestDebounce(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Just("A", "B", "C").Pipe(
			AddLatencyToValues(1, 2),
			operators.Debounce(func(interface{}) rx.Observable {
				return rx.Timer(Step(3))
			}),
		),
		"C", Completed,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			AddLatencyToValues(1, 3),
			operators.Debounce(func(interface{}) rx.Observable {
				return rx.Timer(Step(2))
			}),
		),
		"A", "B", "C", Completed,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			AddLatencyToValues(1, 3),
			operators.Debounce(func(interface{}) rx.Observable {
				return rx.Empty()
			}),
		),
		"C", Completed,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			AddLatencyToValues(1, 3),
			operators.Debounce(func(interface{}) rx.Observable {
				return rx.Throw(ErrTest)
			}),
		),
		ErrTest,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.Debounce(func(interface{}) rx.Observable {
				return rx.Timer(Step(1))
			}),
		),
		ErrTest,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			AddLatencyToValues(1, 2),
			operators.DebounceTime(Step(3)),
		),
		"C", Completed,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			AddLatencyToValues(1, 3),
			operators.DebounceTime(Step(2)),
		),
		"A", "B", "C", Completed,
	).TestAll()
}
