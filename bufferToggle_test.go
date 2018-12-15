package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_BufferToggle(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.BufferToggle(
					Interval(step(2)),
					func(interface{}) Observable { return Interval(step(2)) },
				),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.BufferToggle(
					Interval(step(2)),
					func(interface{}) Observable { return Interval(step(4)) },
				),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.BufferToggle(
					Interval(step(4)),
					func(interface{}) Observable { return Interval(step(2)) },
				),
				toString,
			),
			Concat(Just("A", "B", "C", "D", "E", "F", "G"), Throw(xErrTest)).Pipe(
				addLatencyToNotification(1, 2),
				operators.BufferToggle(
					Interval(step(4)),
					func(interface{}) Observable { return Interval(step(2)) },
				),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.BufferToggle(
					Interval(step(4)),
					func(idx interface{}) Observable {
						if idx.(int) > 1 {
							return Throw(xErrTest)
						}
						return Interval(step(2))
					},
				),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.BufferToggle(
					Interval(step(4)).Pipe(operators.Take(2)),
					func(interface{}) Observable { return Interval(step(2)) },
				),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.BufferToggle(
					Concat(Interval(step(4)).Pipe(operators.Take(2)), Throw(xErrTest)),
					func(interface{}) Observable { return Interval(step(2)) },
				),
				toString,
			),
		},
		"[B]", "[C]", "[D]", "[E]", "[F]", "[G]", xComplete,
		"[B C]", "[C D]", "[D E]", "[E F]", "[F G]", "[G]", xComplete,
		"[C]", "[E]", "[G]", xComplete,
		"[C]", "[E]", "[G]", xErrTest,
		"[C]", "[E]", xErrTest,
		"[C]", "[E]", xComplete,
		"[C]", xErrTest,
	)
}
