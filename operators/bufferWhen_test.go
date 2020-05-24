package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestBufferWhen(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferWhen(func() rx.Observable { return rx.Timer(Step(2)) }),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferWhen(func() rx.Observable { return rx.Timer(Step(4)) }),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferWhen(func() rx.Observable { return rx.Timer(Step(6)) }),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferWhen(func() rx.Observable { return rx.Timer(Step(8)) }),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferWhen(func() rx.Observable { return rx.Throw(ErrTest) }),
				ToString(),
			),
		},
		[][]interface{}{
			{"[A]", "[B]", "[C]", "[D]", "[E]", "[F]", "[G]", rx.Complete},
			{"[A B]", "[C D]", "[E F]", "[G]", rx.Complete},
			{"[A B C]", "[D E F]", "[G]", rx.Complete},
			{"[A B C D]", "[E F G]", rx.Complete},
			{ErrTest},
		},
	)
}
