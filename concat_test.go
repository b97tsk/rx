package rx_test

import (
	"context"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestConcat(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Concat[string](),
		ErrCompleted,
	).Case(
		rx.Concat(
			rx.Pipe1(rx.Just("A", "B"), AddLatencyToValues[string](3, 5)),
			rx.Pipe1(rx.Just("C", "D"), AddLatencyToValues[string](2, 4)),
			rx.Pipe1(rx.Just("E", "F"), AddLatencyToValues[string](1, 3)),
		),
		"A", "B", "C", "D", "E", "F", ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Pipe1(rx.Just("A", "B"), AddLatencyToValues[string](3, 5)),
			rx.ConcatWith(
				rx.Pipe1(rx.Just("C", "D"), AddLatencyToValues[string](2, 4)),
				rx.Pipe1(rx.Just("E", "F"), AddLatencyToValues[string](1, 3)),
			),
		),
		"A", "B", "C", "D", "E", "F", ErrCompleted,
	)

	ctx, cancel := context.WithTimeout(context.Background(), Step(1))
	defer cancel()

	NewTestSuite[int](t).WithContext(ctx).Case(
		rx.Concat(
			func(_ context.Context, sink rx.Observer[int]) {
				time.Sleep(Step(2))
				sink.Complete()
			},
			func(context.Context, rx.Observer[int]) {
				panic("should not happen")
			},
		),
		context.DeadlineExceeded,
	)
}

func TestConcat2(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe1(
			rx.Just(
				rx.Pipe1(rx.Just("A", "B"), AddLatencyToValues[string](3, 5)),
				rx.Pipe1(rx.Just("C", "D"), AddLatencyToValues[string](2, 4)),
				rx.Pipe1(rx.Just("E", "F"), AddLatencyToValues[string](1, 3)),
			),
			rx.ConcatAll[rx.Observable[string]](),
		),
		"A", "B", "C", "D", "E", "F", ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Timer(Step(1)),
			rx.ConcatMap(
				func(time.Time) rx.Observable[string] {
					return rx.Just("A")
				},
			),
		),
		"A", ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Timer(Step(1)),
			rx.ConcatMapTo[time.Time](rx.Just("A")),
		),
		"A", ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Throw[rx.Observable[string]](ErrTest),
			rx.ConcatAll[rx.Observable[string]](),
		),
		ErrTest,
	)

	ctx, cancel := context.WithTimeout(context.Background(), Step(1))
	defer cancel()

	NewTestSuite[string](t).WithContext(ctx).Case(
		rx.Pipe1(
			rx.Just[rx.Observable[string]](
				func(_ context.Context, sink rx.Observer[string]) {
					time.Sleep(Step(2))
					sink.Complete()
				},
				func(context.Context, rx.Observer[string]) {
					panic("should not happen")
				},
			),
			rx.ConcatAll[rx.Observable[string]](),
		),
		context.DeadlineExceeded,
	)
}
