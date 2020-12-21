package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestBufferTime(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferTime(Step(2)),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferTime(Step(4)),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferTime(Step(6)),
				ToString(),
			),
			rx.Throw(ErrTest).Pipe(
				operators.BufferTime(Step(1)),
			),
		},
		[][]interface{}{
			{"[A]", "[B]", "[C]", "[D]", "[E]", Completed},
			{"[A B]", "[C D]", "[E]", Completed},
			{"[A B C]", "[D E]", Completed},
			{ErrTest},
		},
	)

	t.Log("----------")

	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferTimeConfigure{
					TimeSpan: Step(8),
				}.Make(),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferTimeConfigure{
					TimeSpan:      Step(8),
					MaxBufferSize: 3,
				}.Make(),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferTimeConfigure{
					TimeSpan:      Step(8),
					MaxBufferSize: 2,
				}.Make(),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferTimeConfigure{
					TimeSpan:      Step(8),
					MaxBufferSize: 1,
				}.Make(),
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
				operators.BufferTimeConfigure{
					TimeSpan:         Step(2),
					CreationInterval: Step(2),
				}.Make(),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferTimeConfigure{
					TimeSpan:         Step(2),
					CreationInterval: Step(4),
				}.Make(),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferTimeConfigure{
					TimeSpan:         Step(4),
					CreationInterval: Step(2),
				}.Make(),
				ToString(),
			),
		},
		[][]interface{}{
			{"[A]", "[B]", "[C]", "[D]", "[E]", Completed},
			{"[A]", "[C]", "[E]", Completed},
			{"[A B]", "[B C]", "[C D]", "[D E]", "[E]", Completed},
		},
	)
}
