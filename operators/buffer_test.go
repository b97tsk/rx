package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestBuffer(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.Buffer(rx.Ticker(Step(2))),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.Buffer(rx.Ticker(Step(4))),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.Buffer(rx.Ticker(Step(6))),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.Buffer(rx.Ticker(Step(8))),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.Buffer(rx.Throw(ErrTest)),
				ToString(),
			),
			rx.Throw(ErrTest).Pipe(
				operators.Buffer(rx.Throw(ErrTest).Pipe(DelaySubscription(1))),
			),
		},
		[][]interface{}{
			{"[A]", "[B]", "[C]", "[D]", "[E]", "[F]", Completed},
			{"[A B]", "[C D]", "[E F]", Completed},
			{"[A B C]", "[D E F]", Completed},
			{"[A B C D]", Completed},
			{ErrTest},
			{ErrTest},
		},
	)
}
