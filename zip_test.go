package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestZip(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Zip(
			rx.Just("A", "B"),
			rx.Range(1, 4),
		).Pipe(
			ToString(),
		),
		"[A 1]", "[B 2]", Completed,
	).Case(
		rx.Zip(
			rx.Just("A", "B", "C"),
			rx.Range(1, 4),
		).Pipe(
			ToString(),
		),
		"[A 1]", "[B 2]", "[C 3]", Completed,
	).Case(
		rx.Zip(
			rx.Just("A", "B", "C", "D"),
			rx.Range(1, 4),
		).Pipe(
			ToString(),
		),
		"[A 1]", "[B 2]", "[C 3]", Completed,
	).Case(
		rx.Zip(
			rx.Just("A", "B"),
			rx.Concat(
				rx.Range(1, 4),
				rx.Throw(ErrTest),
			).Pipe(
				DelaySubscription(1),
			),
		).Pipe(
			ToString(),
		),
		"[A 1]", "[B 2]", Completed,
	).Case(
		rx.Zip(
			rx.Just("A", "B", "C"),
			rx.Concat(
				rx.Range(1, 4),
				rx.Throw(ErrTest),
			).Pipe(
				DelaySubscription(1),
			),
		).Pipe(
			ToString(),
		),
		"[A 1]", "[B 2]", "[C 3]", Completed,
	).Case(
		rx.Zip(
			rx.Just("A", "B", "C", "D"),
			rx.Concat(
				rx.Range(1, 4),
				rx.Throw(ErrTest),
			).Pipe(
				DelaySubscription(1),
			),
		).Pipe(
			ToString(),
		),
		"[A 1]", "[B 2]", "[C 3]", ErrTest,
	).Case(
		rx.Zip(),
		Completed,
	).TestAll()
}
