package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestFlat(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Just(
			rx.Just("A", "B").Pipe(AddLatencyToValues(3, 5)),
			rx.Just("C", "D").Pipe(AddLatencyToValues(2, 4)),
			rx.Just("E", "F").Pipe(AddLatencyToValues(1, 3)),
		).Pipe(
			operators.Flat(rx.CombineLatest),
			ToString(),
		),
		"[A C E]", "[A C F]", "[A D F]", "[B D F]", Completed,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.Flat(rx.CombineLatest),
		),
		ErrTest,
	).TestAll()

	NewTestSuite(t).Case(
		rx.Just(
			rx.Just("A", "B"),
			rx.Range(1, 4),
		).Pipe(
			operators.Flat(rx.Zip),
			ToString(),
		),
		"[A 1]", "[B 2]", Completed,
	).Case(
		rx.Just(
			rx.Just("A", "B", "C"),
			rx.Range(1, 4),
		).Pipe(
			operators.Flat(rx.Zip),
			ToString(),
		),
		"[A 1]", "[B 2]", "[C 3]", Completed,
	).Case(
		rx.Just(
			rx.Just("A", "B", "C", "D"),
			rx.Range(1, 4),
		).Pipe(
			operators.Flat(rx.Zip),
			ToString(),
		),
		"[A 1]", "[B 2]", "[C 3]", Completed,
	).Case(
		rx.Just(
			rx.Just("A", "B"),
			rx.Concat(rx.Range(1, 4), rx.Throw(ErrTest)).Pipe(DelaySubscription(1)),
		).Pipe(
			operators.Flat(rx.Zip),
			ToString(),
		),
		"[A 1]", "[B 2]", Completed,
	).Case(
		rx.Just(
			rx.Just("A", "B", "C"),
			rx.Concat(rx.Range(1, 4), rx.Throw(ErrTest)).Pipe(DelaySubscription(1)),
		).Pipe(
			operators.Flat(rx.Zip),
			ToString(),
		),
		"[A 1]", "[B 2]", "[C 3]", Completed,
	).Case(
		rx.Just(
			rx.Just("A", "B", "C", "D"),
			rx.Concat(rx.Range(1, 4), rx.Throw(ErrTest)).Pipe(DelaySubscription(1)),
		).Pipe(
			operators.Flat(rx.Zip),
			ToString(),
		),
		"[A 1]", "[B 2]", "[C 3]", ErrTest,
	).TestAll()
}
