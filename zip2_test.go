package rx_test

import (
	"fmt"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestZip2(t *testing.T) {
	t.Parallel()

	toString := func(v1 string, v2 int) string {
		return fmt.Sprintf("[%v %v]", v1, v2)
	}

	NewTestSuite[string](t).Case(
		rx.Zip2(
			rx.Just("A", "B"),
			rx.Pipe1(rx.Range(1, 4), DelaySubscription[int](1)),
			toString,
		),
		"[A 1]", "[B 2]", ErrComplete,
	).Case(
		rx.Zip2(
			rx.Just("A", "B", "C"),
			rx.Pipe1(rx.Range(1, 4), DelaySubscription[int](1)),
			toString,
		),
		"[A 1]", "[B 2]", "[C 3]", ErrComplete,
	).Case(
		rx.Zip2(
			rx.Just("A", "B", "C", "D"),
			rx.Pipe1(rx.Range(1, 4), DelaySubscription[int](1)),
			toString,
		),
		"[A 1]", "[B 2]", "[C 3]", ErrComplete,
	).Case(
		rx.Zip2(
			rx.Just("A", "B"),
			rx.Pipe1(
				rx.Concat(
					rx.Pipe1(rx.Range(1, 4), DelaySubscription[int](1)),
					rx.Throw[int](ErrTest),
				),
				DelaySubscription[int](1),
			),
			toString,
		),
		"[A 1]", "[B 2]", ErrComplete,
	).Case(
		rx.Zip2(
			rx.Just("A", "B", "C"),
			rx.Pipe1(
				rx.Concat(
					rx.Pipe1(rx.Range(1, 4), DelaySubscription[int](1)),
					rx.Throw[int](ErrTest),
				),
				DelaySubscription[int](1),
			),
			toString,
		),
		"[A 1]", "[B 2]", "[C 3]", ErrComplete,
	).Case(
		rx.Zip2(
			rx.Just("A", "B", "C", "D"),
			rx.Pipe1(
				rx.Concat(
					rx.Pipe1(rx.Range(1, 4), DelaySubscription[int](1)),
					rx.Throw[int](ErrTest),
				),
				DelaySubscription[int](1),
			),
			toString,
		),
		"[A 1]", "[B 2]", "[C 3]", ErrTest,
	)
}
