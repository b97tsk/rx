package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/testing"
)

func TestOperators_ForkJoin(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.ForkJoin(
				rx.Just("A", "B", "C").Pipe(AddLatencyToValues(0, 3)),
				rx.Range(1, 5).Pipe(AddLatencyToValues(1, 2)),
				rx.Range(5, 9).Pipe(AddLatencyToValues(3, 1)),
			).Pipe(ToString()),
			rx.ForkJoin(
				rx.Just("A", "B", "C").Pipe(AddLatencyToValues(0, 3)),
				rx.Range(1, 5).Pipe(AddLatencyToValues(1, 2)),
				rx.Range(5, 9).Pipe(AddLatencyToValues(3, 1)),
				rx.Empty().Pipe(DelaySubscription(5)),
			).Pipe(ToString()),
			rx.ForkJoin(
				rx.Just("A", "B", "C").Pipe(AddLatencyToValues(0, 3)),
				rx.Range(1, 5).Pipe(AddLatencyToValues(1, 2)),
				rx.Range(5, 9).Pipe(AddLatencyToValues(3, 1)),
				rx.Throw(ErrTest).Pipe(DelaySubscription(5)),
			).Pipe(ToString()),
		},
		[][]interface{}{
			{"[C 4 8]", rx.Completed},
			{rx.Completed},
			{ErrTest},
		},
	)
}
