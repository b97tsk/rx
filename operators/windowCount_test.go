package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestWindowCount(t *testing.T) {
	toSlice := func(val interface{}, idx int) rx.Observable {
		if obs, ok := val.(rx.Observable); ok {
			return obs.Pipe(operators.ToSlice())
		}

		return rx.Throw(rx.ErrNotObservable)
	}

	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowCount(2),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowCount(3),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowCountConfigure{
					WindowSize:       3,
					StartWindowEvery: 1,
				}.Make(),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowCountConfigure{
					WindowSize:       3,
					StartWindowEvery: 2,
				}.Make(),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowCountConfigure{
					WindowSize:       3,
					StartWindowEvery: 4,
				}.Make(),
				operators.MergeMap(toSlice),
				ToString(),
			),
		},
		[][]interface{}{
			{"[A B]", "[C D]", "[E F]", "[G]", Completed},
			{"[A B C]", "[D E F]", "[G]", Completed},
			{"[A B C]", "[B C D]", "[C D E]", "[D E F]", "[E F G]", "[F G]", "[G]", "[]", Completed},
			{"[A B C]", "[C D E]", "[E F G]", "[G]", Completed},
			{"[A B C]", "[E F G]", "[]", Completed},
		},
	)
}
