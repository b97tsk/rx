package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestThrottle(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(0, 2),
			operators.Throttle(func(interface{}) rx.Observable {
				return rx.Timer(Step(3))
			}),
		),
		"A", "C", "E", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(0, 2),
			operators.Throttle(func(interface{}) rx.Observable {
				return rx.Empty()
			}),
		),
		"A", "B", "C", "D", "E", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(0, 2),
			operators.ThrottleConfigure{
				DurationSelector: func(interface{}) rx.Observable {
					return rx.Empty().Pipe(DelaySubscription(5))
				},
				Leading:  false,
				Trailing: true,
			}.Make(),
		),
		Completed,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.Throttle(func(interface{}) rx.Observable {
				return rx.Throw(ErrTest)
			}),
		),
		ErrTest,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(0, 2),
			operators.Throttle(func(interface{}) rx.Observable {
				return rx.Throw(ErrTest)
			}),
		),
		"A", ErrTest,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(0, 4),
			operators.ThrottleConfigure{
				DurationSelector: func(interface{}) rx.Observable {
					return rx.Timer(Step(9))
				},
				Leading:  false,
				Trailing: true,
			}.Make(),
		),
		"C", "E", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(0, 4),
			operators.ThrottleConfigure{
				DurationSelector: func(interface{}) rx.Observable {
					return rx.Timer(Step(9))
				},
				Leading:  true,
				Trailing: true,
			}.Make(),
		),
		"A", "C", "E", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(0, 2),
			operators.ThrottleTime(Step(3)),
		),
		"A", "C", "E", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(0, 4),
			operators.ThrottleTimeConfigure{
				Duration: Step(9),
				Leading:  false,
				Trailing: true,
			}.Make(),
		),
		"C", "E", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(0, 4),
			operators.ThrottleTimeConfigure{
				Duration: Step(9),
				Leading:  true,
				Trailing: true,
			}.Make(),
		),
		"A", "C", "E", Completed,
	).TestAll()

	panictest := func(f func(), msg string) {
		defer func() {
			if recover() == nil {
				t.Log(msg)
				t.FailNow()
			}
		}()
		f()
	}
	panictest(
		func() { operators.Throttle(nil) },
		"Throttle with nil duration selector didn't panic.",
	)
}
