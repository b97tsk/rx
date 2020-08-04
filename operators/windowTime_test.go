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
				operators.WindowTimeConfigure{
					TimeSpan: Step(8),
				}.Make(),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTimeConfigure{
					TimeSpan:      Step(8),
					MaxWindowSize: 3,
				}.Make(),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTimeConfigure{
					TimeSpan:      Step(8),
					MaxWindowSize: 2,
				}.Make(),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTimeConfigure{
					TimeSpan:      Step(8),
					MaxWindowSize: 1,
				}.Make(),
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
				operators.WindowTimeConfigure{
					TimeSpan:         Step(2),
					CreationInterval: Step(2),
				}.Make(),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTimeConfigure{
					TimeSpan:         Step(2),
					CreationInterval: Step(4),
				}.Make(),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTimeConfigure{
					TimeSpan:         Step(4),
					CreationInterval: Step(2),
				}.Make(),
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
