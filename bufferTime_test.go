package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_BufferTime(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.BufferTime(step(2), 0),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.BufferTime(step(4), 0),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.BufferTime(step(6), 0),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.BufferTime(step(8), 0),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.BufferTime(step(8), 3),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.BufferTime(step(8), 2),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.BufferTime(step(8), 1),
				toString,
			),
		},
		"[A]", "[B]", "[C]", "[D]", "[E]", "[F]", "[G]", xComplete,
		"[A B]", "[C D]", "[E F]", "[G]", xComplete,
		"[A B C]", "[D E F]", "[G]", xComplete,
		"[A B C D]", "[E F G]", xComplete,
		"[A B C]", "[D E F]", "[G]", xComplete,
		"[A B]", "[C D]", "[E F]", "[G]", xComplete,
		"[A]", "[B]", "[C]", "[D]", "[E]", "[F]", "[G]", "[]", xComplete,
	)
	t.Log("----------")
	subscribe(
		t,
		[]Observable{
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				BufferTimeConfigure{step(2), step(2), 0}.MakeFunc(),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				BufferTimeConfigure{step(2), step(4), 0}.MakeFunc(),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				BufferTimeConfigure{step(4), step(2), 0}.MakeFunc(),
				toString,
			),
		},
		"[A]", "[B]", "[C]", "[D]", "[E]", "[F]", "[G]", xComplete,
		"[A]", "[C]", "[E]", "[G]", xComplete,
		"[A B]", "[B C]", "[C D]", "[D E]", "[E F]", "[F G]", "[G]", xComplete,
	)
}
