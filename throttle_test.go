package rx_test

import (
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestThrottle(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			AddLatencyToValues[string](0, 2),
			rx.Throttle(
				func(string) rx.Observable[time.Time] {
					return rx.Timer(Step(3))
				},
			).AsOperator(),
		),
		"A", "C", "E", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			AddLatencyToValues[string](0, 2),
			rx.Throttle(
				func(string) rx.Observable[int] {
					return rx.Empty[int]()
				},
			).AsOperator(),
		),
		"A", "B", "C", "D", "E", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			AddLatencyToValues[string](0, 2),
			rx.Throttle(
				func(string) rx.Observable[int] {
					return rx.Pipe1(
						rx.Empty[int](),
						DelaySubscription[int](5),
					)
				},
			).WithLeading(false).WithTrailing(true).AsOperator(),
		),
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Throw[string](ErrTest),
			rx.Throttle(
				func(string) rx.Observable[int] {
					return rx.Throw[int](ErrTest)
				},
			).AsOperator(),
		),
		ErrTest,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			AddLatencyToValues[string](0, 2),
			rx.Throttle(
				func(string) rx.Observable[int] {
					return rx.Throw[int](ErrTest)
				},
			).AsOperator(),
		),
		"A", ErrTest,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			AddLatencyToValues[string](0, 4),
			rx.Throttle(
				func(string) rx.Observable[time.Time] {
					return rx.Timer(Step(9))
				},
			).WithLeading(false).WithTrailing(true).AsOperator(),
		),
		"C", "E", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			AddLatencyToValues[string](0, 4),
			rx.Throttle(
				func(string) rx.Observable[time.Time] {
					return rx.Timer(Step(9))
				},
			).WithLeading(true).WithTrailing(true).AsOperator(),
		),
		"A", "C", "E", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			AddLatencyToValues[string](0, 2),
			rx.ThrottleTime[string](Step(3)).AsOperator(),
		),
		"A", "C", "E", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			AddLatencyToValues[string](0, 4),
			rx.ThrottleTime[string](Step(9)).WithLeading(false).WithTrailing(true).AsOperator(),
		),
		"C", "E", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			AddLatencyToValues[string](0, 4),
			rx.ThrottleTime[string](Step(9)).WithLeading(true).WithTrailing(true).AsOperator(),
		),
		"A", "C", "E", ErrComplete,
	)
}
