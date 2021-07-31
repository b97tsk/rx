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
			operators.BufferCountConfig{
				BufferSize:       3,
				StartBufferEvery: 1,
			}.Make(),
			ToString(),
		),
		"[A B C]", "[B C D]", "[C D E]", "[D E]", "[E]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			operators.BufferCountConfig{
				BufferSize:       3,
				StartBufferEvery: 2,
			}.Make(),
			ToString(),
		),
		"[A B C]", "[C D E]", "[E]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			operators.BufferCountConfig{
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

	panictest := func(f func(), msg string) {
		defer func() {
			if recover() == nil {
				t.Log(msg)
				t.FailNow()
			}
		}()
		f()
	}
	panictest(
		func() { operators.BufferCount(-1) },
		"BufferCount with negative buffer size didn't panic.",
	)
	panictest(
		func() { operators.BufferCount(0) },
		"BufferCount with zero buffer size didn't panic.",
	)
	panictest(
		func() {
			operators.BufferCountConfig{
				BufferSize:       1,
				StartBufferEvery: -1,
			}.Make()
		},
		"BufferCountConfig with negative StartBufferEvery didn't panic.",
	)
}
