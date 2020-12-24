package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestBufferWhen(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferWhen(func() rx.Observable { return rx.Timer(Step(2)) }),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferWhen(func() rx.Observable { return rx.Timer(Step(4)) }),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferWhen(func() rx.Observable { return rx.Timer(Step(6)) }),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferWhen(func() rx.Observable { return rx.Timer(Step(8)) }),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferWhen(rx.Empty),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferWhen(func() rx.Observable { return rx.Throw(ErrTest) }),
				ToString(),
			),
			rx.Throw(ErrTest).Pipe(
				operators.BufferWhen(func() rx.Observable { return rx.Timer(Step(1)) }),
			),
		},
		[][]interface{}{
			{"[A]", "[B]", "[C]", "[D]", "[E]", Completed},
			{"[A B]", "[C D]", "[E]", Completed},
			{"[A B C]", "[D E]", Completed},
			{"[A B C D]", "[E]", Completed},
			{"[A B C D E]", Completed},
			{ErrTest},
			{ErrTest},
		},
	)
}
