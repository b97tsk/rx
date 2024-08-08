package rx_test

import (
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
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Range(1, 10),
			rx.TakeLast[int](3),
		),
		7, 8, 9, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Just(1),
			rx.TakeLast[int](3),
		),
		1, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Empty[int](),
			rx.TakeLast[int](3),
		),
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 10),
				rx.Throw[int](ErrTest),
			),
			rx.TakeLast[int](0),
		),
		ErrComplete,
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
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 10),
				rx.Oops[int](ErrTest),
			),
			rx.TakeLast[int](0),
		),
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 10),
				rx.Oops[int](ErrTest),
			),
			rx.TakeLast[int](3),
		),
		rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just(1),
				rx.Oops[int](ErrTest),
			),
			rx.TakeLast[int](3),
		),
		rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Empty[int](),
				rx.Oops[int](ErrTest),
			),
			rx.TakeLast[int](3),
		),
		rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe2(
			rx.Range(1, 10),
			rx.TakeLast[int](3),
			rx.DoOnNext(func(int) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)

	ctx, cancel := rx.NewBackgroundContext().WithTimeoutCause(Step(1), ErrTest)
	defer cancel()

	NewTestSuite[int](t).WithContext(ctx).Case(
		rx.Pipe2(
			rx.Range(1, 10),
			rx.TakeLast[int](3),
			rx.DoOnNext(func(int) { time.Sleep(Step(2)) }),
		),
		7, ErrTest,
	)
}
