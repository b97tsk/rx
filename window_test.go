package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Window(t *testing.T) {
	toSlice := func(val interface{}, idx int) Observable {
		if obs, ok := val.(Observable); ok {
			return obs.Pipe(
				operators.ToSlice(),
			)
		}
		return Throw(ErrNotObservable)
	}
	subscribe(
		t,
		[]Observable{
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.Window(Interval(step(2))),
				operators.MergeMap(toSlice),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.Window(Interval(step(4))),
				operators.MergeMap(toSlice),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.Window(Interval(step(6))),
				operators.MergeMap(toSlice),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.Window(Interval(step(8))),
				operators.MergeMap(toSlice),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.Window(Throw(xErrTest)),
				operators.MergeMap(toSlice),
				toString,
			),
		},
		"[A]", "[B]", "[C]", "[D]", "[E]", "[F]", "[G]", xComplete,
		"[A B]", "[C D]", "[E F]", "[G]", xComplete,
		"[A B C]", "[D E F]", "[G]", xComplete,
		"[A B C D]", "[E F G]", xComplete,
		xErrTest,
	)
}
