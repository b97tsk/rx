package rx_test

import (
	"fmt"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestZipWithBuffering3(t *testing.T) {
	t.Parallel()

	proj := func(v1, v2 string, v3 int) string {
		return fmt.Sprintf("[%v %v %v]", v1, v2, v3)
	}

	NewTestSuite[string](t).Case(
		rx.ZipWithBuffering3(
			rx.Just("A", "B"),
			rx.Just("B", "C"),
			rx.Pipe1(rx.Range(1, 4), DelaySubscription[int](1)),
			proj,
		),
		"[A B 1]", "[B C 2]", ErrComplete,
	).Case(
		rx.ZipWithBuffering3(
			rx.Just("A", "B", "C"),
			rx.Just("B", "C", "D"),
			rx.Pipe1(rx.Range(1, 4), DelaySubscription[int](1)),
			proj,
		),
		"[A B 1]", "[B C 2]", "[C D 3]", ErrComplete,
	).Case(
		rx.ZipWithBuffering3(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Pipe1(rx.Range(1, 4), DelaySubscription[int](1)),
			proj,
		),
		"[A B 1]", "[B C 2]", "[C D 3]", ErrComplete,
	).Case(
		rx.ZipWithBuffering3(
			rx.Just("A", "B"),
			rx.Just("B", "C"),
			rx.Pipe1(
				rx.Concat(
					rx.Pipe1(rx.Range(1, 4), DelaySubscription[int](1)),
					rx.Throw[int](ErrTest),
				),
				DelaySubscription[int](1),
			),
			proj,
		),
		"[A B 1]", "[B C 2]", ErrComplete,
	).Case(
		rx.ZipWithBuffering3(
			rx.Just("A", "B", "C"),
			rx.Just("B", "C", "D"),
			rx.Pipe1(
				rx.Concat(
					rx.Pipe1(rx.Range(1, 4), DelaySubscription[int](1)),
					rx.Throw[int](ErrTest),
				),
				DelaySubscription[int](1),
			),
			proj,
		),
		"[A B 1]", "[B C 2]", "[C D 3]", ErrComplete,
	).Case(
		rx.ZipWithBuffering3(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Pipe1(
				rx.Concat(
					rx.Pipe1(rx.Range(1, 4), DelaySubscription[int](1)),
					rx.Throw[int](ErrTest),
				),
				DelaySubscription[int](1),
			),
			proj,
		),
		"[A B 1]", "[B C 2]", "[C D 3]", ErrTest,
	).Case(
		rx.Pipe1(
			rx.ZipWithBuffering3(
				rx.Just("A", "B", "C", "D"),
				rx.Just("B", "C", "D", "E"),
				rx.Range(1, 5),
				proj,
			),
			rx.OnNext(func(string) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}
