package rx_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestZip7(t *testing.T) {
	t.Parallel()

	toString := func(v1, v2, v3, v4, v5, v6 string, v7 int) string {
		return fmt.Sprintf("[%v %v %v %v %v %v %v]", v1, v2, v3, v4, v5, v6, v7)
	}

	NewTestSuite[string](t).Case(
		rx.Zip7(
			rx.Just("A", "B"),
			rx.Just("B", "C"),
			rx.Just("C", "D"),
			rx.Just("D", "E"),
			rx.Just("E", "F"),
			rx.Just("F", "G"),
			rx.Pipe(rx.Range(1, 4), DelaySubscription[int](1)),
			toString,
		),
		"[A B C D E F 1]", "[B C D E F G 2]", ErrCompleted,
	).Case(
		rx.Zip7(
			rx.Just("A", "B", "C"),
			rx.Just("B", "C", "D"),
			rx.Just("C", "D", "E"),
			rx.Just("D", "E", "F"),
			rx.Just("E", "F", "G"),
			rx.Just("F", "G", "H"),
			rx.Pipe(rx.Range(1, 4), DelaySubscription[int](1)),
			toString,
		),
		"[A B C D E F 1]", "[B C D E F G 2]", "[C D E F G H 3]", ErrCompleted,
	).Case(
		rx.Zip7(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			rx.Just("D", "E", "F", "G"),
			rx.Just("E", "F", "G", "H"),
			rx.Just("F", "G", "H", "I"),
			rx.Pipe(rx.Range(1, 4), DelaySubscription[int](1)),
			toString,
		),
		"[A B C D E F 1]", "[B C D E F G 2]", "[C D E F G H 3]", ErrCompleted,
	).Case(
		rx.Zip7(
			rx.Just("A", "B"),
			rx.Just("B", "C"),
			rx.Just("C", "D"),
			rx.Just("D", "E"),
			rx.Just("E", "F"),
			rx.Just("F", "G"),
			rx.Pipe(
				rx.Concat(
					rx.Pipe(rx.Range(1, 4), DelaySubscription[int](1)),
					rx.Throw[int](ErrTest),
				),
				DelaySubscription[int](1),
			),
			toString,
		),
		"[A B C D E F 1]", "[B C D E F G 2]", ErrCompleted,
	).Case(
		rx.Zip7(
			rx.Just("A", "B", "C"),
			rx.Just("B", "C", "D"),
			rx.Just("C", "D", "E"),
			rx.Just("D", "E", "F"),
			rx.Just("E", "F", "G"),
			rx.Just("F", "G", "H"),
			rx.Pipe(
				rx.Concat(
					rx.Pipe(rx.Range(1, 4), DelaySubscription[int](1)),
					rx.Throw[int](ErrTest),
				),
				DelaySubscription[int](1),
			),
			toString,
		),
		"[A B C D E F 1]", "[B C D E F G 2]", "[C D E F G H 3]", ErrCompleted,
	).Case(
		rx.Zip7(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			rx.Just("D", "E", "F", "G"),
			rx.Just("E", "F", "G", "H"),
			rx.Just("F", "G", "H", "I"),
			rx.Pipe(
				rx.Concat(
					rx.Pipe(rx.Range(1, 4), DelaySubscription[int](1)),
					rx.Throw[int](ErrTest),
				),
				DelaySubscription[int](1),
			),
			toString,
		),
		"[A B C D E F 1]", "[B C D E F G 2]", "[C D E F G H 3]", ErrTest,
	)

	ctx, cancel := context.WithTimeout(context.Background(), Step(1))
	defer cancel()

	NewTestSuite[string](t).WithContext(ctx).Case(
		rx.Zip7(
			rx.Just("A"),
			rx.Just("B"),
			rx.Just("C"),
			rx.Just("D"),
			rx.Just("E"),
			rx.Just("F"),
			rx.Timer(Step(2)),
			func(v1, v2, v3, v4, v5, v6 string, _ time.Time) string {
				return v1 + v2 + v3 + v4 + v5 + v6
			},
		),
		context.DeadlineExceeded,
	)
}