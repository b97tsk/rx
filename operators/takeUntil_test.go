package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestTakeUntil(t *testing.T) {
	addLatency := AddLatencyToValues(0, 2)
	delay := DelaySubscription(3)
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C").Pipe(addLatency, operators.TakeUntil(rx.Just(42))),
			rx.Just("A", "B", "C").Pipe(addLatency, operators.TakeUntil(rx.Empty())),
			rx.Just("A", "B", "C").Pipe(addLatency, operators.TakeUntil(rx.Throw(ErrTest))),
			rx.Just("A", "B", "C").Pipe(addLatency, operators.TakeUntil(rx.Never())),
		},
		[][]interface{}{
			{Completed},
			{Completed},
			{ErrTest},
			{"A", "B", "C", Completed},
		},
	)
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C").Pipe(addLatency, operators.TakeUntil(rx.Just(42).Pipe(delay))),
			rx.Just("A", "B", "C").Pipe(addLatency, operators.TakeUntil(rx.Empty().Pipe(delay))),
			rx.Just("A", "B", "C").Pipe(addLatency, operators.TakeUntil(rx.Throw(ErrTest).Pipe(delay))),
		},
		[][]interface{}{
			{"A", "B", Completed},
			{"A", "B", Completed},
			{"A", "B", ErrTest},
		},
	)
}
