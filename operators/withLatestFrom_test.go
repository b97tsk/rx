package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestWithLatestFrom(t *testing.T) {
	addLatency1 := AddLatencyToValues(1, 2)
	addLatency2 := AddLatencyToNotifications(0, 2)

	observables := [...]rx.Observable{
		rx.Just("A", "B").Pipe(addLatency1),
		rx.Just("A", "B", "C").Pipe(addLatency1),
		rx.Just("A", "B", "C", "D").Pipe(addLatency1),
	}

	{
		observables := observables
		for i, obs := range observables {
			observables[i] = obs.Pipe(
				operators.WithLatestFrom(rx.Range(1, 4).Pipe(addLatency2)),
				ToString(),
			)
		}
		SubscribeN(
			t,
			observables[:],
			[][]interface{}{
				{"[A 1]", "[B 2]", Completed},
				{"[A 1]", "[B 2]", "[C 3]", Completed},
				{"[A 1]", "[B 2]", "[C 3]", "[D 3]", Completed},
			},
		)
	}

	{
		observables := observables
		for i, obs := range observables {
			observables[i] = obs.Pipe(
				operators.WithLatestFrom(rx.Concat(rx.Range(1, 4), rx.Throw(ErrTest)).Pipe(addLatency2)),
				ToString(),
			)
		}
		SubscribeN(
			t,
			observables[:],
			[][]interface{}{
				{"[A 1]", "[B 2]", Completed},
				{"[A 1]", "[B 2]", "[C 3]", Completed},
				{"[A 1]", "[B 2]", "[C 3]", ErrTest},
			},
		)
	}
}
