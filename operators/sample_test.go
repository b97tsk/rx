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
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.Sample(rx.Ticker(Step(4))),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.Sample(
					rx.Concat(
						rx.Ticker(Step(4)).Pipe(operators.Take(3)),
						rx.Throw(ErrTest),
					),
				),
			),
		},
		[][]interface{}{
			{"B", "D", "F", rx.Completed},
			{"B", "D", "F", ErrTest},
		},
	)
}

func TestSampleTime(t *testing.T) {
	Subscribe(
		t,
		rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
			AddLatencyToValues(1, 2),
			operators.SampleTime(Step(4)),
		),
		"B", "D", "F", rx.Completed,
	)
}
