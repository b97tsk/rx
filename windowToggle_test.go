package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_WindowToggle(t *testing.T) {
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
				operators.WindowToggle(
					Interval(step(2)),
					func(interface{}) Observable { return Interval(step(2)) },
				),
				operators.MergeMap(toSlice),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.WindowToggle(
					Interval(step(2)),
					func(interface{}) Observable { return Interval(step(4)) },
				),
				operators.MergeMap(toSlice),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.WindowToggle(
					Interval(step(4)),
					func(interface{}) Observable { return Interval(step(2)) },
				),
				operators.MergeMap(toSlice),
				toString,
			),
			Concat(Just("A", "B", "C", "D", "E", "F", "G"), Throw(errTest)).Pipe(
				addLatencyToNotification(1, 2),
				operators.WindowToggle(
					Interval(step(4)),
					func(interface{}) Observable { return Interval(step(2)) },
				),
				operators.MergeMap(toSlice),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.WindowToggle(
					Interval(step(4)),
					func(idx interface{}) Observable {
						if idx.(int) > 1 {
							return Throw(errTest)
						}
						return Interval(step(2))
					},
				),
				operators.MergeMap(toSlice),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.WindowToggle(
					Interval(step(4)).Pipe(operators.Take(2)),
					func(interface{}) Observable { return Interval(step(2)) },
				),
				operators.MergeMap(toSlice),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.WindowToggle(
					Concat(Interval(step(4)).Pipe(operators.Take(2)), Throw(errTest)),
					func(interface{}) Observable { return Interval(step(2)) },
				),
				operators.MergeMap(toSlice),
				toString,
			),
		},
		"[B]", "[C]", "[D]", "[E]", "[F]", "[G]", Complete,
		"[B C]", "[C D]", "[D E]", "[E F]", "[F G]", "[G]", Complete,
		"[C]", "[E]", "[G]", Complete,
		"[C]", "[E]", "[G]", errTest,
		"[C]", "[E]", errTest,
		"[C]", "[E]", Complete,
		"[C]", errTest,
	)
}
