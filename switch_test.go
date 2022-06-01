package rx_test

import (
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestSwitch(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe2(
			rx.Just(
				rx.Pipe(rx.Just("A", "B", "C", "D"), AddLatencyToValues[string](0, 2)),
				rx.Pipe(rx.Just("E", "F", "G", "H"), AddLatencyToValues[string](0, 3)),
				rx.Pipe(rx.Just("I", "J", "K", "L"), AddLatencyToValues[string](0, 2)),
			),
			AddLatencyToValues[rx.Observable[string]](0, 5),
			rx.SwitchAll[rx.Observable[string]](),
		),
		"A", "B", "C", "E", "F", "I", "J", "K", "L", ErrCompleted,
	).Case(
		rx.Pipe2(
			rx.Just(
				rx.Pipe(rx.Just("A", "B", "C", "D"), AddLatencyToValues[string](0, 2)),
				rx.Pipe(rx.Just("E", "F", "G", "H"), AddLatencyToValues[string](0, 3)),
				rx.Pipe(rx.Just("I", "J", "K", "L"), AddLatencyToValues[string](0, 2)),
				rx.Empty[string](),
			),
			AddLatencyToValues[rx.Observable[string]](0, 5),
			rx.SwitchAll[rx.Observable[string]](),
		),
		"A", "B", "C", "E", "F", "I", "J", "K", ErrCompleted,
	).Case(
		rx.Pipe2(
			rx.Just(
				rx.Pipe(rx.Just("A", "B", "C", "D"), AddLatencyToValues[string](0, 2)),
				rx.Pipe(rx.Just("E", "F", "G", "H"), AddLatencyToValues[string](0, 3)),
				rx.Pipe(rx.Just("I", "J", "K", "L"), AddLatencyToValues[string](0, 2)),
				rx.Throw[string](ErrTest),
			),
			AddLatencyToValues[rx.Observable[string]](0, 5),
			rx.SwitchAll[rx.Observable[string]](),
		),
		"A", "B", "C", "E", "F", "I", "J", "K", ErrTest,
	).Case(
		rx.Pipe(
			rx.Timer(Step(1)),
			rx.SwitchMap(
				func(time.Time) rx.Observable[string] {
					return rx.Just("A")
				},
			),
		),
		"A", ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Timer(Step(1)),
			rx.SwitchMapTo[time.Time](rx.Just("A")),
		),
		"A", ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Throw[rx.Observable[string]](ErrTest),
			rx.SwitchAll[rx.Observable[string]](),
		),
		ErrTest,
	)
}
