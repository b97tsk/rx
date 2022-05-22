package rx_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestZip3(t *testing.T) {
	t.Parallel()

	toString := func(v1, v2 string, v3 int) string {
		return fmt.Sprintf("[%v %v %v]", v1, v2, v3)
	}

	NewTestSuite[string](t).Case(
		rx.Zip3(
			rx.Just("A", "B"),
			rx.Just("B", "C"),
			rx.Pipe(rx.Range(1, 4), DelaySubscription[int](1)),
			toString,
		),
		"[A B 1]", "[B C 2]", ErrCompleted,
	).Case(
		rx.Zip3(
			rx.Just("A", "B", "C"),
			rx.Just("B", "C", "D"),
			rx.Pipe(rx.Range(1, 4), DelaySubscription[int](1)),
			toString,
		),
		"[A B 1]", "[B C 2]", "[C D 3]", ErrCompleted,
	).Case(
		rx.Zip3(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Pipe(rx.Range(1, 4), DelaySubscription[int](1)),
			toString,
		),
		"[A B 1]", "[B C 2]", "[C D 3]", ErrCompleted,
	).Case(
		rx.Zip3(
			rx.Just("A", "B"),
			rx.Just("B", "C"),
			rx.Pipe(
				rx.Concat(
					rx.Pipe(rx.Range(1, 4), DelaySubscription[int](1)),
					rx.Throw[int](ErrTest),
				),
				DelaySubscription[int](1),
			),
			toString,
		),
		"[A B 1]", "[B C 2]", ErrCompleted,
	).Case(
		rx.Zip3(
			rx.Just("A", "B", "C"),
			rx.Just("B", "C", "D"),
			rx.Pipe(
				rx.Concat(
					rx.Pipe(rx.Range(1, 4), DelaySubscription[int](1)),
					rx.Throw[int](ErrTest),
				),
				DelaySubscription[int](1),
			),
			toString,
		),
		"[A B 1]", "[B C 2]", "[C D 3]", ErrCompleted,
	).Case(
		rx.Zip3(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Pipe(
				rx.Concat(
					rx.Pipe(rx.Range(1, 4), DelaySubscription[int](1)),
					rx.Throw[int](ErrTest),
				),
				DelaySubscription[int](1),
			),
			toString,
		),
		"[A B 1]", "[B C 2]", "[C D 3]", ErrTest,
	)

	ctx, cancel := context.WithTimeout(context.Background(), Step(1))
	defer cancel()

	NewTestSuite[string](t).WithContext(ctx).Case(
		rx.Zip3(
			rx.Just("A"),
			rx.Just("B"),
			rx.Timer(Step(2)),
			func(v1, v2 string, _ time.Time) string {
				return v1 + v2
			},
		),
		context.DeadlineExceeded,
	)
}
