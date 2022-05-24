package rx_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestZip5(t *testing.T) {
	t.Parallel()

	toString := func(v1, v2, v3, v4 string, v5 int) string {
		return fmt.Sprintf("[%v %v %v %v %v]", v1, v2, v3, v4, v5)
	}

	NewTestSuite[string](t).Case(
		rx.Zip5(
			rx.Just("A", "B"),
			rx.Just("B", "C"),
			rx.Just("C", "D"),
			rx.Just("D", "E"),
			rx.Pipe(rx.Range(1, 4), DelaySubscription[int](1)),
			toString,
		),
		"[A B C D 1]", "[B C D E 2]", ErrCompleted,
	).Case(
		rx.Zip5(
			rx.Just("A", "B", "C"),
			rx.Just("B", "C", "D"),
			rx.Just("C", "D", "E"),
			rx.Just("D", "E", "F"),
			rx.Pipe(rx.Range(1, 4), DelaySubscription[int](1)),
			toString,
		),
		"[A B C D 1]", "[B C D E 2]", "[C D E F 3]", ErrCompleted,
	).Case(
		rx.Zip5(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			rx.Just("D", "E", "F", "G"),
			rx.Pipe(rx.Range(1, 4), DelaySubscription[int](1)),
			toString,
		),
		"[A B C D 1]", "[B C D E 2]", "[C D E F 3]", ErrCompleted,
	).Case(
		rx.Zip5(
			rx.Just("A", "B"),
			rx.Just("B", "C"),
			rx.Just("C", "D"),
			rx.Just("D", "E"),
			rx.Pipe(
				rx.Concat(
					rx.Pipe(rx.Range(1, 4), DelaySubscription[int](1)),
					rx.Throw[int](ErrTest),
				),
				DelaySubscription[int](1),
			),
			toString,
		),
		"[A B C D 1]", "[B C D E 2]", ErrCompleted,
	).Case(
		rx.Zip5(
			rx.Just("A", "B", "C"),
			rx.Just("B", "C", "D"),
			rx.Just("C", "D", "E"),
			rx.Just("D", "E", "F"),
			rx.Pipe(
				rx.Concat(
					rx.Pipe(rx.Range(1, 4), DelaySubscription[int](1)),
					rx.Throw[int](ErrTest),
				),
				DelaySubscription[int](1),
			),
			toString,
		),
		"[A B C D 1]", "[B C D E 2]", "[C D E F 3]", ErrCompleted,
	).Case(
		rx.Zip5(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			rx.Just("D", "E", "F", "G"),
			rx.Pipe(
				rx.Concat(
					rx.Pipe(rx.Range(1, 4), DelaySubscription[int](1)),
					rx.Throw[int](ErrTest),
				),
				DelaySubscription[int](1),
			),
			toString,
		),
		"[A B C D 1]", "[B C D E 2]", "[C D E F 3]", ErrTest,
	)

	ctx, cancel := context.WithTimeout(context.Background(), Step(1))
	defer cancel()

	NewTestSuite[string](t).WithContext(ctx).Case(
		rx.Zip5(
			rx.Just("A"),
			rx.Just("B"),
			rx.Just("C"),
			rx.Just("D"),
			rx.Timer(Step(2)),
			func(v1, v2, v3, v4 string, _ time.Time) string {
				return v1 + v2 + v3 + v4
			},
		),
		context.DeadlineExceeded,
	)
}