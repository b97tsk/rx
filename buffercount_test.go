package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestBufferCount(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			rx.BufferCount[string](2),
			ToString[[]string](),
		),
		"[A B]", "[C D]", "[E]", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			rx.BufferCount[string](3),
			ToString[[]string](),
		),
		"[A B C]", "[D E]", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			rx.BufferCount[string](3).WithStartBufferEvery(1),
			ToString[[]string](),
		),
		"[A B C]", "[B C D]", "[C D E]", "[D E]", "[E]", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			rx.BufferCount[string](3).WithStartBufferEvery(2),
			ToString[[]string](),
		),
		"[A B C]", "[C D E]", "[E]", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			rx.BufferCount[string](3).WithStartBufferEvery(4),
			ToString[[]string](),
		),
		"[A B C]", "[E]", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Throw[string](ErrTest),
			rx.BufferCount[string](2),
			ToString[[]string](),
		),
		ErrTest,
	).Case(
		rx.Pipe2(
			rx.Oops[string](ErrTest),
			rx.BufferCount[string](2),
			ToString[[]string](),
		),
		rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe2(
			rx.Empty[string](),
			rx.BufferCount[string](0),
			ToString[[]string](),
		),
		rx.ErrOops, "BufferCount: BufferSize <= 0",
	).Case(
		rx.Pipe2(
			rx.Empty[string](),
			rx.BufferCount[string](2).WithStartBufferEvery(0),
			ToString[[]string](),
		),
		rx.ErrOops, "BufferCount: StartBufferEvery <= 0",
	).Case(
		rx.Pipe3(
			rx.Just("A"),
			rx.BufferCount[string](2),
			rx.DoOnNext(func([]string) { panic(ErrTest) }),
			ToString[[]string](),
		),
		rx.ErrOops, ErrTest,
	)
}
