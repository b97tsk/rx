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
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTime(Step(2)),
				operators.MergeMap(toSlice, -1),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTime(Step(4)),
				operators.MergeMap(toSlice, -1),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTime(Step(6)),
				operators.MergeMap(toSlice, -1),
				ToString(),
			),
		},
		[][]interface{}{
			{"[A]", "[B]", "[C]", "[D]", "[E]", Completed},
			{"[A B]", "[C D]", "[E]", Completed},
			{"[A B C]", "[D E]", Completed},
		},
	)

	t.Log("----------")

	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTimeConfigure{
					TimeSpan: Step(8),
				}.Make(),
				operators.MergeMap(toSlice, -1),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTimeConfigure{
					TimeSpan:      Step(8),
					MaxWindowSize: 3,
				}.Make(),
				operators.MergeMap(toSlice, -1),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTimeConfigure{
					TimeSpan:      Step(8),
					MaxWindowSize: 2,
				}.Make(),
				operators.MergeMap(toSlice, -1),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTimeConfigure{
					TimeSpan:      Step(8),
					MaxWindowSize: 1,
				}.Make(),
				operators.MergeMap(toSlice, -1),
				ToString(),
			),
		},
		[][]interface{}{
			{"[A B C D]", "[E]", Completed},
			{"[A B C]", "[D E]", Completed},
			{"[A B]", "[C D]", "[E]", Completed},
			{"[A]", "[B]", "[C]", "[D]", "[E]", "[]", Completed},
		},
	)

	t.Log("----------")

	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTimeConfigure{
					TimeSpan:         Step(2),
					CreationInterval: Step(2),
				}.Make(),
				operators.MergeMap(toSlice, -1),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTimeConfigure{
					TimeSpan:         Step(2),
					CreationInterval: Step(4),
				}.Make(),
				operators.MergeMap(toSlice, -1),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTimeConfigure{
					TimeSpan:         Step(4),
					CreationInterval: Step(2),
				}.Make(),
				operators.MergeMap(toSlice, -1),
				ToString(),
			),
			rx.Throw(ErrTest).Pipe(
				operators.WindowTime(Step(1)),
				operators.MergeMap(toSlice, -1),
				ToString(),
			),
		},
		[][]interface{}{
			{"[A]", "[B]", "[C]", "[D]", "[E]", Completed},
			{"[A]", "[C]", "[E]", Completed},
			{"[A B]", "[B C]", "[C D]", "[D E]", "[E]", Completed},
			{ErrTest},
		},
	)
}
