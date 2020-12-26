package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestBuffer(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
			AddLatencyToValues(1, 2),
			operators.Buffer(rx.Ticker(Step(2))),
			ToString(),
		),
		"[A]", "[B]", "[C]", "[D]", "[E]", "[F]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
			AddLatencyToValues(1, 2),
			operators.Buffer(rx.Ticker(Step(4))),
			ToString(),
		),
		"[A B]", "[C D]", "[E F]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
			AddLatencyToValues(1, 2),
			operators.Buffer(rx.Ticker(Step(6))),
			ToString(),
		),
		"[A B C]", "[D E F]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
			AddLatencyToValues(1, 2),
			operators.Buffer(rx.Ticker(Step(8))),
			ToString(),
		),
		"[A B C D]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
			AddLatencyToValues(1, 2),
			operators.Buffer(rx.Throw(ErrTest)),
			ToString(),
		),
		ErrTest,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.Buffer(rx.Throw(ErrTest).Pipe(DelaySubscription(1))),
		),
		ErrTest,
	).TestAll()
}
