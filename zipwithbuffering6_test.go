package rx_test

import (
	"fmt"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestZipWithBuffering6(t *testing.T) {
	t.Parallel()

	toString := func(v1, v2, v3, v4, v5 string, v6 int) string {
		return fmt.Sprintf("[%v %v %v %v %v %v]", v1, v2, v3, v4, v5, v6)
	}

	NewTestSuite[string](t).Case(
		rx.ZipWithBuffering6(
			rx.Just("A", "B"),
			rx.Just("B", "C"),
			rx.Just("C", "D"),
			rx.Just("D", "E"),
			rx.Just("E", "F"),
			rx.Pipe1(rx.Range(1, 4), DelaySubscription[int](1)),
			toString,
		),
		"[A B C D E 1]", "[B C D E F 2]", ErrComplete,
	).Case(
		rx.ZipWithBuffering6(
			rx.Just("A", "B", "C"),
			rx.Just("B", "C", "D"),
			rx.Just("C", "D", "E"),
			rx.Just("D", "E", "F"),
			rx.Just("E", "F", "G"),
			rx.Pipe1(rx.Range(1, 4), DelaySubscription[int](1)),
			toString,
		),
		"[A B C D E 1]", "[B C D E F 2]", "[C D E F G 3]", ErrComplete,
	).Case(
		rx.ZipWithBuffering6(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			rx.Just("D", "E", "F", "G"),
			rx.Just("E", "F", "G", "H"),
			rx.Pipe1(rx.Range(1, 4), DelaySubscription[int](1)),
			toString,
		),
		"[A B C D E 1]", "[B C D E F 2]", "[C D E F G 3]", ErrComplete,
	).Case(
		rx.ZipWithBuffering6(
			rx.Just("A", "B"),
			rx.Just("B", "C"),
			rx.Just("C", "D"),
			rx.Just("D", "E"),
			rx.Just("E", "F"),
			rx.Pipe1(
				rx.Concat(
					rx.Pipe1(rx.Range(1, 4), DelaySubscription[int](1)),
					rx.Throw[int](ErrTest),
				),
				DelaySubscription[int](1),
			),
			toString,
		),
		"[A B C D E 1]", "[B C D E F 2]", ErrComplete,
	).Case(
		rx.ZipWithBuffering6(
			rx.Just("A", "B", "C"),
			rx.Just("B", "C", "D"),
			rx.Just("C", "D", "E"),
			rx.Just("D", "E", "F"),
			rx.Just("E", "F", "G"),
			rx.Pipe1(
				rx.Concat(
					rx.Pipe1(rx.Range(1, 4), DelaySubscription[int](1)),
					rx.Throw[int](ErrTest),
				),
				DelaySubscription[int](1),
			),
			toString,
		),
		"[A B C D E 1]", "[B C D E F 2]", "[C D E F G 3]", ErrComplete,
	).Case(
		rx.ZipWithBuffering6(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			rx.Just("D", "E", "F", "G"),
			rx.Just("E", "F", "G", "H"),
			rx.Pipe1(
				rx.Concat(
					rx.Pipe1(rx.Range(1, 4), DelaySubscription[int](1)),
					rx.Throw[int](ErrTest),
				),
				DelaySubscription[int](1),
			),
			toString,
		),
		"[A B C D E 1]", "[B C D E F 2]", "[C D E F G 3]", ErrTest,
	)
}
