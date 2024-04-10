package rx_test

import (
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestExhaust(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe2(
			rx.Just(
				rx.Pipe1(rx.Just("A", "B", "C", "D"), AddLatencyToValues[string](0, 2)),
				rx.Pipe1(rx.Just("E", "F", "G", "H"), AddLatencyToValues[string](0, 3)),
				rx.Pipe1(rx.Just("I", "J", "K", "L"), AddLatencyToValues[string](0, 2)),
			),
			AddLatencyToValues[rx.Observable[string]](0, 5),
			rx.ExhaustAll[rx.Observable[string]](),
		),
		"A", "B", "C", "D", "I", "J", "K", "L", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just(
				rx.Pipe1(rx.Just("A", "B", "C", "D"), AddLatencyToValues[string](0, 2)),
				rx.Pipe1(rx.Just("E", "F", "G", "H"), AddLatencyToValues[string](0, 3)),
				rx.Pipe1(rx.Just("I", "J", "K", "L"), AddLatencyToValues[string](0, 2)),
				rx.Throw[string](ErrTest),
			),
			AddLatencyToValues[rx.Observable[string]](0, 5),
			rx.ExhaustAll[rx.Observable[string]](),
		),
		"A", "B", "C", "D", "I", "J", "K", "L", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just(
				rx.Pipe1(rx.Just("A", "B", "C", "D"), AddLatencyToValues[string](0, 2)),
				rx.Pipe1(rx.Just("E", "F", "G", "H"), AddLatencyToValues[string](0, 3)),
				rx.Pipe1(rx.Just("I", "J", "K", "L"), AddLatencyToValues[string](0, 2)),
				rx.Throw[string](ErrTest),
				rx.Throw[string](ErrTest),
			),
			AddLatencyToValues[rx.Observable[string]](0, 5),
			rx.ExhaustAll[rx.Observable[string]](),
		),
		"A", "B", "C", "D", "I", "J", "K", "L", ErrTest,
	).Case(
		rx.Pipe1(
			rx.Timer(Step(1)),
			rx.ExhaustMap(
				func(time.Time) rx.Observable[string] {
					return rx.Just("A")
				},
			),
		),
		"A", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Timer(Step(1)),
			rx.ExhaustMapTo[time.Time](rx.Just("A")),
		),
		"A", ErrComplete,
	).Case(
		rx.Pipe1(
			func(_ rx.Context, o rx.Observer[rx.Observable[string]]) {
				o.Next(rx.Pipe1(rx.Just("A", "B", "C", "D"), AddLatencyToValues[string](0, 2)))
				time.Sleep(Step(5))
				o.Error(ErrTest)
			},
			rx.ExhaustAll[rx.Observable[string]](),
		),
		"A", "B", "C", ErrTest,
	).Case(
		rx.Pipe1(
			rx.Timer(Step(1)),
			rx.ExhaustMap(func(time.Time) rx.Observable[string] { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}
