package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestBufferWhen(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.BufferWhen(func() rx.Observable { return rx.Timer(Step(2)) }),
			ToString(),
		),
		"[A]", "[B]", "[C]", "[D]", "[E]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.BufferWhen(func() rx.Observable { return rx.Timer(Step(4)) }),
			ToString(),
		),
		"[A B]", "[C D]", "[E]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.BufferWhen(func() rx.Observable { return rx.Timer(Step(6)) }),
			ToString(),
		),
		"[A B C]", "[D E]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.BufferWhen(func() rx.Observable { return rx.Timer(Step(8)) }),
			ToString(),
		),
		"[A B C D]", "[E]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.BufferWhen(rx.Empty),
			ToString(),
		),
		"[A B C D E]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.BufferWhen(func() rx.Observable { return rx.Throw(ErrTest) }),
			ToString(),
		),
		ErrTest,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.BufferWhen(func() rx.Observable { return rx.Timer(Step(1)) }),
		),
		ErrTest,
	).TestAll()
}
