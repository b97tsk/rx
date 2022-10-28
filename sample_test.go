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
		"B", "D", ErrCompleted,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			AddLatencyToValues[string](1, 2),
			rx.Sample[string](rx.Empty[int]()),
		),
		ErrCompleted,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			AddLatencyToValues[string](1, 2),
			rx.Sample[string](rx.Never[int]()),
		),
		ErrCompleted,
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
					rx.Pipe(
						rx.Ticker(Step(4)),
						rx.Take[time.Time](2),
					),
					rx.Throw[time.Time](ErrTest),
				),
			),
		),
		"B", "D", ErrTest,
	).Case(
		rx.Pipe(
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
		"B", "D", ErrCompleted,
	)

	ctx, cancel := context.WithTimeout(context.Background(), Step(1))
	defer cancel()

	NewTestSuite[string](t).WithContext(ctx).Case(
		rx.Pipe(
			rx.Empty[string](),
			rx.Sample[string](
				func(ctx context.Context, sink rx.Observer[string]) {
					time.Sleep(Step(2))
					sink.Complete()
				},
			),
		),
		context.DeadlineExceeded,
	)
}