package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestBufferTime(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferTime(Step(2)),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferTime(Step(4)),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferTime(Step(6)),
				ToString(),
			),
		},
		[][]interface{}{
			{"[A]", "[B]", "[C]", "[D]", "[E]", "[F]", "[G]", rx.Complete},
			{"[A B]", "[C D]", "[E F]", "[G]", rx.Complete},
			{"[A B C]", "[D E F]", "[G]", rx.Complete},
		},
	)
	t.Log("----------")
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferTimeConfigure{Step(8), 0, 0}.Use(),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferTimeConfigure{Step(8), 0, 3}.Use(),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferTimeConfigure{Step(8), 0, 2}.Use(),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferTimeConfigure{Step(8), 0, 1}.Use(),
				ToString(),
			),
		},
		[][]interface{}{
			{"[A B C D]", "[E F G]", rx.Complete},
			{"[A B C]", "[D E F]", "[G]", rx.Complete},
			{"[A B]", "[C D]", "[E F]", "[G]", rx.Complete},
			{"[A]", "[B]", "[C]", "[D]", "[E]", "[F]", "[G]", "[]", rx.Complete},
		},
	)
	t.Log("----------")
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferTimeConfigure{Step(2), Step(2), 0}.Use(),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferTimeConfigure{Step(2), Step(4), 0}.Use(),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferTimeConfigure{Step(4), Step(2), 0}.Use(),
				ToString(),
			),
		},
		[][]interface{}{
			{"[A]", "[B]", "[C]", "[D]", "[E]", "[F]", "[G]", rx.Complete},
			{"[A]", "[C]", "[E]", "[G]", rx.Complete},
			{"[A B]", "[B C]", "[C D]", "[D E]", "[E F]", "[F G]", "[G]", rx.Complete},
		},
	)
}
