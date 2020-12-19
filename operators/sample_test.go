package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestSample(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.Sample(rx.Ticker(Step(4))),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.Sample(rx.Empty()),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.Sample(
					rx.Concat(
						rx.Ticker(Step(4)).Pipe(operators.Take(2)),
						rx.Throw(ErrTest),
					),
				),
			),
			rx.Throw(ErrTest).Pipe(
				operators.Sample(rx.Ticker(Step(1))),
			),
		},
		[][]interface{}{
			{"B", "D", Completed},
			{Completed},
			{"B", "D", ErrTest},
			{ErrTest},
		},
	)
}

func TestSampleTime(t *testing.T) {
	Subscribe(
		t,
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.SampleTime(Step(4)),
		),
		"B", "D", Completed,
	)
}
