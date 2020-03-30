package operators_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_WindowTime(t *testing.T) {
	toSlice := func(val interface{}, idx int) Observable {
		if obs, ok := val.(Observable); ok {
			return obs.Pipe(
				operators.ToSlice(),
			)
		}
		return Throw(ErrNotObservable)
	}
	subscribeN(
		t,
		[]Observable{
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.WindowTime(step(2)),
				operators.MergeMap(toSlice),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.WindowTime(step(4)),
				operators.MergeMap(toSlice),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.WindowTime(step(6)),
				operators.MergeMap(toSlice),
				toString,
			),
		},
		[][]interface{}{
			{"[A]", "[B]", "[C]", "[D]", "[E]", "[F]", "[G]", Complete},
			{"[A B]", "[C D]", "[E F]", "[G]", Complete},
			{"[A B C]", "[D E F]", "[G]", Complete},
		},
	)
	t.Log("----------")
	subscribeN(
		t,
		[]Observable{
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				WindowTimeConfigure{step(8), 0, 0}.Use(),
				operators.MergeMap(toSlice),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				WindowTimeConfigure{step(8), 0, 3}.Use(),
				operators.MergeMap(toSlice),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				WindowTimeConfigure{step(8), 0, 2}.Use(),
				operators.MergeMap(toSlice),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				WindowTimeConfigure{step(8), 0, 1}.Use(),
				operators.MergeMap(toSlice),
				toString,
			),
		},
		[][]interface{}{
			{"[A B C D]", "[E F G]", Complete},
			{"[A B C]", "[D E F]", "[G]", Complete},
			{"[A B]", "[C D]", "[E F]", "[G]", Complete},
			{"[A]", "[B]", "[C]", "[D]", "[E]", "[F]", "[G]", "[]", Complete},
		},
	)
	t.Log("----------")
	subscribeN(
		t,
		[]Observable{
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				WindowTimeConfigure{step(2), step(2), 0}.Use(),
				operators.MergeMap(toSlice),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				WindowTimeConfigure{step(2), step(4), 0}.Use(),
				operators.MergeMap(toSlice),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				WindowTimeConfigure{step(4), step(2), 0}.Use(),
				operators.MergeMap(toSlice),
				toString,
			),
		},
		[][]interface{}{
			{"[A]", "[B]", "[C]", "[D]", "[E]", "[F]", "[G]", Complete},
			{"[A]", "[C]", "[E]", "[G]", Complete},
			{"[A B]", "[B C]", "[C D]", "[D E]", "[E F]", "[F G]", "[G]", Complete},
		},
	)
}
