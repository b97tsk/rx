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
			Concat(Just("A", "B", "C", "D", "E", "F", "G"), Throw(errTest)).Pipe(
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
							return Throw(errTest)
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
					Concat(Interval(step(4)).Pipe(operators.Take(2)), Throw(errTest)),
					func(interface{}) Observable { return Interval(step(2)) },
				),
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
