package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestInterval(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.CombineLatest(
				rx.Interval(Step(1)).Pipe(operators.Take(3)),
				rx.Interval(0).Pipe(operators.Take(3)),
			).Pipe(ToString()),
			rx.CombineLatest(
				rx.Interval(Step(2)).Pipe(operators.Take(3)),
				rx.Timer(Step(1), Step(2)).Pipe(operators.Take(3)),
			).Pipe(ToString()),
		},
		[][]interface{}{
			{"[0 2]", "[1 2]", "[2 2]", rx.Complete},
			{"[0 0]", "[0 1]", "[1 1]", "[1 2]", "[2 2]", rx.Complete},
		},
	)
}
