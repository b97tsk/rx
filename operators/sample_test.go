package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestSample(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.Sample(rx.Interval(Step(4))),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.Sample(
					rx.Concat(
						rx.Interval(Step(4)).Pipe(operators.Take(3)),
						rx.Throw(ErrTest),
					),
				),
			),
		},
		[][]interface{}{
			{"B", "D", "F", rx.Complete},
			{"B", "D", "F", ErrTest},
		},
	)
}