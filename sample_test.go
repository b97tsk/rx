package rx_test

import (
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
			rx.SampleTime[string](Step(4)),
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
			rx.Empty[string](),
			rx.Sample[string](rx.Throw[int](ErrTest)),
		),
		ErrComplete,
	)
}
