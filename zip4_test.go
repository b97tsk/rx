package rx_test

import (
	"fmt"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestZip4(t *testing.T) {
	t.Parallel()

	toString := func(v1, v2, v3 string, v4 int) string {
		return fmt.Sprintf("[%v %v %v %v]", v1, v2, v3, v4)
	}

	NewTestSuite[string](t).Case(
		rx.Zip4(
			rx.Just("A", "B"),
			rx.Just("B", "C"),
			rx.Just("C", "D"),
			rx.Pipe(rx.Range(1, 4), DelaySubscription[int](1)),
			toString,
		),
		"[A B C 1]", "[B C D 2]", ErrCompleted,
	).Case(
		rx.Zip4(
			rx.Just("A", "B", "C"),
			rx.Just("B", "C", "D"),
			rx.Just("C", "D", "E"),
			rx.Pipe(rx.Range(1, 4), DelaySubscription[int](1)),
			toString,
		),
		"[A B C 1]", "[B C D 2]", "[C D E 3]", ErrCompleted,
	).Case(
		rx.Zip4(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			rx.Pipe(rx.Range(1, 4), DelaySubscription[int](1)),
			toString,
		),
		"[A B C 1]", "[B C D 2]", "[C D E 3]", ErrCompleted,
	).Case(
		rx.Zip4(
			rx.Just("A", "B"),
			rx.Just("B", "C"),
			rx.Just("C", "D"),
			rx.Pipe(
				rx.Concat(
					rx.Pipe(rx.Range(1, 4), DelaySubscription[int](1)),
					rx.Throw[int](ErrTest),
				),
				DelaySubscription[int](1),
			),
			toString,
		),
		"[A B C 1]", "[B C D 2]", ErrCompleted,
	).Case(
		rx.Zip4(
			rx.Just("A", "B", "C"),
			rx.Just("B", "C", "D"),
			rx.Just("C", "D", "E"),
			rx.Pipe(
				rx.Concat(
					rx.Pipe(rx.Range(1, 4), DelaySubscription[int](1)),
					rx.Throw[int](ErrTest),
				),
				DelaySubscription[int](1),
			),
			toString,
		),
		"[A B C 1]", "[B C D 2]", "[C D E 3]", ErrCompleted,
	).Case(
		rx.Zip4(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			rx.Pipe(
				rx.Concat(
					rx.Pipe(rx.Range(1, 4), DelaySubscription[int](1)),
					rx.Throw[int](ErrTest),
				),
				DelaySubscription[int](1),
			),
			toString,
		),
		"[A B C 1]", "[B C D 2]", "[C D E 3]", ErrTest,
	)
}
