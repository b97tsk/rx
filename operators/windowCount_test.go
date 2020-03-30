package operators_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_WindowCount(t *testing.T) {
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
				operators.WindowCount(2),
				operators.MergeMap(toSlice),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.WindowCount(3),
				operators.MergeMap(toSlice),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				WindowCountConfigure{3, 1}.Use(),
				operators.MergeMap(toSlice),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				WindowCountConfigure{3, 2}.Use(),
				operators.MergeMap(toSlice),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				WindowCountConfigure{3, 4}.Use(),
				operators.MergeMap(toSlice),
				toString,
			),
		},
		[][]interface{}{
			{"[A B]", "[C D]", "[E F]", "[G]", Complete},
			{"[A B C]", "[D E F]", "[G]", Complete},
			{"[A B C]", "[B C D]", "[C D E]", "[D E F]", "[E F G]", "[F G]", "[G]", "[]", Complete},
			{"[A B C]", "[C D E]", "[E F G]", "[G]", Complete},
			{"[A B C]", "[E F G]", "[]", Complete},
		},
	)
}
