package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestSkipUntil(t *testing.T) {
	addLatency := AddLatencyToValues(0, 2)
	delay := DelaySubscription(3)
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C").Pipe(addLatency, operators.SkipUntil(rx.Just(42))),
			rx.Just("A", "B", "C").Pipe(addLatency, operators.SkipUntil(rx.Empty())),
			rx.Just("A", "B", "C").Pipe(addLatency, operators.SkipUntil(rx.Throw(ErrTest))),
			rx.Just("A", "B", "C").Pipe(addLatency, operators.SkipUntil(rx.Never())),
		},
		[][]interface{}{
			{"A", "B", "C", rx.Completed},
			{rx.Completed},
			{ErrTest},
			{rx.Completed},
		},
	)
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C").Pipe(addLatency, operators.SkipUntil(rx.Just(42).Pipe(delay))),
			rx.Just("A", "B", "C").Pipe(addLatency, operators.SkipUntil(rx.Empty().Pipe(delay))),
			rx.Just("A", "B", "C").Pipe(addLatency, operators.SkipUntil(rx.Throw(ErrTest).Pipe(delay))),
		},
		[][]interface{}{
			{"C", rx.Completed},
			{rx.Completed},
			{ErrTest},
		},
	)
}
