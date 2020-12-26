package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestBufferCount(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			operators.BufferCount(2),
			ToString(),
		),
		"[A B]", "[C D]", "[E]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			operators.BufferCount(3),
			ToString(),
		),
		"[A B C]", "[D E]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			operators.BufferCountConfigure{
				BufferSize:       3,
				StartBufferEvery: 1,
			}.Make(),
			ToString(),
		),
		"[A B C]", "[B C D]", "[C D E]", "[D E]", "[E]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			operators.BufferCountConfigure{
				BufferSize:       3,
				StartBufferEvery: 2,
			}.Make(),
			ToString(),
		),
		"[A B C]", "[C D E]", "[E]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			operators.BufferCountConfigure{
				BufferSize:       3,
				StartBufferEvery: 4,
			}.Make(),
			ToString(),
		),
		"[A B C]", "[E]", Completed,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.BufferCount(2),
		),
		ErrTest,
	).TestAll()
}
