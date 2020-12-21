package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestBufferCount(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C", "D", "E").Pipe(operators.BufferCount(2), ToString()),
			rx.Just("A", "B", "C", "D", "E").Pipe(operators.BufferCount(3), ToString()),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				operators.BufferCountConfigure{
					BufferSize:       3,
					StartBufferEvery: 1,
				}.Make(),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				operators.BufferCountConfigure{
					BufferSize:       3,
					StartBufferEvery: 2,
				}.Make(),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				operators.BufferCountConfigure{
					BufferSize:       3,
					StartBufferEvery: 4,
				}.Make(),
				ToString(),
			),
			rx.Throw(ErrTest).Pipe(operators.BufferCount(2)),
		},
		[][]interface{}{
			{"[A B]", "[C D]", "[E]", Completed},
			{"[A B C]", "[D E]", Completed},
			{"[A B C]", "[B C D]", "[C D E]", "[D E]", "[E]", Completed},
			{"[A B C]", "[C D E]", "[E]", Completed},
			{"[A B C]", "[E]", Completed},
			{ErrTest},
		},
	)
}
