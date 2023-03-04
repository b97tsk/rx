package rx_test

import (
	"context"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestTakeLast(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe1(
			rx.Range(1, 10),
			rx.TakeLast[int](0),
		),
		ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Range(1, 10),
			rx.TakeLast[int](3),
		),
		7, 8, 9, ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Just(1),
			rx.TakeLast[int](3),
		),
		1, ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Empty[int](),
			rx.TakeLast[int](3),
		),
		ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 10),
				rx.Throw[int](ErrTest),
			),
			rx.TakeLast[int](0),
		),
		ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 10),
				rx.Throw[int](ErrTest),
			),
			rx.TakeLast[int](3),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just(1),
				rx.Throw[int](ErrTest),
			),
			rx.TakeLast[int](3),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Empty[int](),
				rx.Throw[int](ErrTest),
			),
			rx.TakeLast[int](3),
		),
		ErrTest,
	)

	ctx, cancel := context.WithTimeout(context.Background(), Step(1))
	defer cancel()

	NewTestSuite[int](t).WithContext(ctx).Case(
		rx.Pipe2(
			rx.Range(1, 10),
			rx.TakeLast[int](3),
			rx.OnNext(func(int) { time.Sleep(Step(2)) }),
		),
		7, context.DeadlineExceeded,
	)
}
