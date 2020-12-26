package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestSample(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.Sample(rx.Ticker(Step(4))),
		),
		"B", "D", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.Sample(rx.Empty()),
		),
		Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.Sample(
				rx.Concat(
					rx.Ticker(Step(4)).Pipe(operators.Take(2)),
					rx.Throw(ErrTest),
				),
			),
		),
		"B", "D", ErrTest,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.Sample(rx.Ticker(Step(1))),
		),
		ErrTest,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.SampleTime(Step(4)),
		),
		"B", "D", Completed,
	).TestAll()
}
