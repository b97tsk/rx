package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestFlat(t *testing.T) {
	t.Run("CombineLatest", func(t *testing.T) {
		Subscribe(
			t,
			rx.Just(
				rx.Just("A", "B").Pipe(AddLatencyToValues(3, 5)),
				rx.Just("C", "D").Pipe(AddLatencyToValues(2, 4)),
				rx.Just("E", "F").Pipe(AddLatencyToValues(1, 3)),
			).Pipe(operators.Flat(rx.CombineLatest), ToString()),
			"[A C E]", "[A C F]", "[A D F]", "[B D F]", rx.Completed,
		)
	})
	t.Run("Zip", func(t *testing.T) {
		delay := DelaySubscription(1)
		observables := [...]rx.Observable{
			rx.Just(rx.Just("A", "B"), rx.Range(1, 4)),
			rx.Just(rx.Just("A", "B", "C"), rx.Range(1, 4)),
			rx.Just(rx.Just("A", "B", "C", "D"), rx.Range(1, 4)),
			rx.Just(rx.Just("A", "B"), rx.Concat(rx.Range(1, 4), rx.Throw(ErrTest)).Pipe(delay)),
			rx.Just(rx.Just("A", "B", "C"), rx.Concat(rx.Range(1, 4), rx.Throw(ErrTest)).Pipe(delay)),
			rx.Just(rx.Just("A", "B", "C", "D"), rx.Concat(rx.Range(1, 4), rx.Throw(ErrTest)).Pipe(delay)),
		}
		for i, obs := range observables {
			observables[i] = obs.Pipe(operators.Flat(rx.Zip), ToString())
		}
		SubscribeN(
			t,
			observables[:],
			[][]interface{}{
				{"[A 1]", "[B 2]", rx.Completed},
				{"[A 1]", "[B 2]", "[C 3]", rx.Completed},
				{"[A 1]", "[B 2]", "[C 3]", rx.Completed},
				{"[A 1]", "[B 2]", rx.Completed},
				{"[A 1]", "[B 2]", "[C 3]", rx.Completed},
				{"[A 1]", "[B 2]", "[C 3]", ErrTest},
			},
		)
	})
}
