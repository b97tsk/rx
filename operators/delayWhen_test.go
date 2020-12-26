package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestDelayWhen(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			operators.DelayWhen(
				func(val interface{}, idx int) rx.Observable {
					return rx.Timer(Step(idx + 1))
				},
			),
		),
		"A", "B", "C", "D", "E", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToNotifications(0, 2),
			operators.DelayWhen(
				func(interface{}, int) rx.Observable {
					return rx.Timer(Step(1))
				},
			),
		),
		"A", "B", "C", "D", "E", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			operators.DelayWhen(
				func(interface{}, int) rx.Observable {
					return rx.Empty().Pipe(DelaySubscription(1))
				},
			),
		),
		Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			operators.DelayWhen(
				func(interface{}, int) rx.Observable {
					return rx.Throw(ErrTest)
				},
			),
		),
		ErrTest,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.DelayWhen(
				func(interface{}, int) rx.Observable {
					return rx.Timer(Step(1))
				},
			),
		),
		ErrTest,
	).TestAll()
}
