package rx_test

import (
	"context"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestSample(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			AddLatencyToValues[string](1, 2),
			rx.Sample[string](rx.Ticker(Step(4))),
		),
		"B", "D", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			AddLatencyToValues[string](1, 2),
			rx.Sample[string](rx.Empty[int]()),
		),
		ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			AddLatencyToValues[string](1, 2),
			rx.Sample[string](rx.Never[int]()),
		),
		ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			AddLatencyToValues[string](1, 2),
			rx.Sample[string](rx.Throw[int](ErrTest)),
		),
		ErrTest,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			AddLatencyToValues[string](1, 2),
			rx.Sample[string](
				rx.Concat(
					rx.Pipe1(
						rx.Ticker(Step(4)),
						rx.Take[time.Time](2),
					),
					rx.Throw[time.Time](ErrTest),
				),
			),
		),
		"B", "D", ErrTest,
	).Case(
		rx.Pipe1(
			rx.Throw[string](ErrTest),
			rx.Sample[string](rx.Ticker(Step(1))),
		),
		ErrTest,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			AddLatencyToValues[string](1, 2),
			rx.SampleTime[string](Step(4)),
		),
		"B", "D", ErrComplete,
	)

	ctx, cancel := rx.NewBackgroundContext().WithTimeout(Step(1))
	defer cancel()

	NewTestSuite[string](t).WithContext(ctx).Case(
		rx.Pipe1(
			rx.Empty[string](),
			rx.Sample[string](
				func(_ rx.Context, sink rx.Observer[string]) {
					time.Sleep(Step(2))
					sink.Complete()
				},
			),
		),
		context.DeadlineExceeded,
	)

	t.Run("Oops", func(t *testing.T) {
		defer func() {
			NewTestSuite[string](t).Case(rx.Oops[string](recover()), rx.ErrOops, ErrTest)
		}()
		rx.Pipe1(
			rx.Empty[string](),
			rx.Sample[string](
				func(_ rx.Context, sink rx.Observer[int]) {
					defer sink.Complete()
					panic(ErrTest)
				},
			),
		).Subscribe(rx.NewBackgroundContext(), rx.Noop[string])
	})
}
