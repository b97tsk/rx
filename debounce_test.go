package rx_test

import (
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestDebounce(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](1, 2),
			rx.Debounce(
				func(string) rx.Observable[time.Time] {
					return rx.Timer(Step(3))
				},
			),
		),
		"C", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](1, 3),
			rx.Debounce(
				func(string) rx.Observable[time.Time] {
					return rx.Timer(Step(2))
				},
			),
		),
		"A", "B", "C", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](1, 3),
			rx.Debounce(
				func(string) rx.Observable[int] {
					return rx.Empty[int]()
				},
			),
		),
		"C", ErrComplete,
	).Case(
		rx.Pipe1(
			func(_ rx.Context, sink rx.Observer[string]) {
				sink.Next("A")
				time.Sleep(Step(1))
				sink.Next("B")
				time.Sleep(Step(1))
				sink.Next("C")
				sink.Complete()
			},
			rx.Debounce(
				func(string) rx.Observable[int] {
					return rx.Throw[int](ErrTest)
				},
			),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.Throw[string](ErrTest),
			rx.Debounce(
				func(string) rx.Observable[time.Time] {
					return rx.Timer(Step(1))
				},
			),
		),
		ErrTest,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](1, 2),
			rx.DebounceTime[string](Step(3)),
		),
		"C", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](1, 3),
			rx.DebounceTime[string](Step(2)),
		),
		"A", "B", "C", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Just("A", "B", "C"),
			rx.Debounce(func(string) rx.Observable[int] { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe3(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](1, 2),
			rx.DebounceTime[string](Step(3)),
			rx.OnNext(func(string) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe3(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](1, 3),
			rx.DebounceTime[string](Step(2)),
			rx.OnNext(func(string) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}
