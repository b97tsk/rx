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
			rx.BufferCount[string](2).AsOperator(),
			ToString[[]string](),
		),
		"[A B]", "[C D]", "[E]", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			rx.BufferCount[string](3).AsOperator(),
			ToString[[]string](),
		),
		"[A B C]", "[D E]", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			rx.BufferCount[string](3).WithStartBufferEvery(1).AsOperator(),
			ToString[[]string](),
		),
		"[A B C]", "[B C D]", "[C D E]", "[D E]", "[E]", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			rx.BufferCount[string](3).WithStartBufferEvery(2).AsOperator(),
			ToString[[]string](),
		),
		"[A B C]", "[C D E]", "[E]", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			rx.BufferCount[string](3).WithStartBufferEvery(4).AsOperator(),
			ToString[[]string](),
		),
		"[A B C]", "[E]", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Throw[string](ErrTest),
			rx.BufferCount[string](2).AsOperator(),
			ToString[[]string](),
		),
		ErrTest,
	)
}
