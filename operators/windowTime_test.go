package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestWindowTime(t *testing.T) {
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
				operators.WindowTime(Step(2)),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTime(Step(4)),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTime(Step(6)),
				operators.MergeMap(toSlice),
				ToString(),
			),
		},
		[][]interface{}{
			{"[A]", "[B]", "[C]", "[D]", "[E]", "[F]", "[G]", Completed},
			{"[A B]", "[C D]", "[E F]", "[G]", Completed},
			{"[A B C]", "[D E F]", "[G]", Completed},
		},
	)
	t.Log("----------")
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTimeConfigure{Step(8), 0, 0}.Make(),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTimeConfigure{Step(8), 0, 3}.Make(),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTimeConfigure{Step(8), 0, 2}.Make(),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTimeConfigure{Step(8), 0, 1}.Make(),
				operators.MergeMap(toSlice),
				ToString(),
			),
		},
		[][]interface{}{
			{"[A B C D]", "[E F G]", Completed},
			{"[A B C]", "[D E F]", "[G]", Completed},
			{"[A B]", "[C D]", "[E F]", "[G]", Completed},
			{"[A]", "[B]", "[C]", "[D]", "[E]", "[F]", "[G]", "[]", Completed},
		},
	)
	t.Log("----------")
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTimeConfigure{Step(2), Step(2), 0}.Make(),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTimeConfigure{Step(2), Step(4), 0}.Make(),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTimeConfigure{Step(4), Step(2), 0}.Make(),
				operators.MergeMap(toSlice),
				ToString(),
			),
		},
		[][]interface{}{
			{"[A]", "[B]", "[C]", "[D]", "[E]", "[F]", "[G]", Completed},
			{"[A]", "[C]", "[E]", "[G]", Completed},
			{"[A B]", "[B C]", "[C D]", "[D E]", "[E F]", "[F G]", "[G]", Completed},
		},
	)
}
