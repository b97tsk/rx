package operators_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_BufferTime(t *testing.T) {
	subscribeN(
		t,
		[]Observable{
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.BufferTime(step(2)),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.BufferTime(step(4)),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.BufferTime(step(6)),
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
				BufferTimeConfigure{step(8), 0, 0}.Use(),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				BufferTimeConfigure{step(8), 0, 3}.Use(),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				BufferTimeConfigure{step(8), 0, 2}.Use(),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				BufferTimeConfigure{step(8), 0, 1}.Use(),
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
				BufferTimeConfigure{step(2), step(2), 0}.Use(),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				BufferTimeConfigure{step(2), step(4), 0}.Use(),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				BufferTimeConfigure{step(4), step(2), 0}.Use(),
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
