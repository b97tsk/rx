package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_ForkJoin(t *testing.T) {
	subscribeN(
		t,
		[]Observable{
			ForkJoin(
				Just("A", "B", "C").Pipe(addLatencyToValue(0, 3)),
				Range(1, 5).Pipe(addLatencyToValue(1, 2)),
				Range(5, 9).Pipe(addLatencyToValue(3, 1)),
			).Pipe(toString),
			ForkJoin(
				Just("A", "B", "C").Pipe(addLatencyToValue(0, 3)),
				Range(1, 5).Pipe(addLatencyToValue(1, 2)),
				Range(5, 9).Pipe(addLatencyToValue(3, 1)),
				Empty().Pipe(delaySubscription(5)),
			).Pipe(toString),
			ForkJoin(
				Just("A", "B", "C").Pipe(addLatencyToValue(0, 3)),
				Range(1, 5).Pipe(addLatencyToValue(1, 2)),
				Range(5, 9).Pipe(addLatencyToValue(3, 1)),
				Throw(errTest).Pipe(delaySubscription(5)),
			).Pipe(toString),
		},
		[][]interface{}{
			{"[C 4 8]", Complete},
			{Complete},
			{errTest},
		},
	)
}
