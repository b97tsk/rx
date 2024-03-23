package rx_test

import (
	"fmt"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestZipWithBuffering7(t *testing.T) {
	t.Parallel()

	mapping := func(v1, v2, v3, v4, v5, v6 string, v7 int) string {
		return fmt.Sprintf("[%v %v %v %v %v %v %v]", v1, v2, v3, v4, v5, v6, v7)
	}

	NewTestSuite[string](t).Case(
		rx.ZipWithBuffering7(
			rx.Just("A", "B"),
			rx.Just("B", "C"),
			rx.Just("C", "D"),
			rx.Just("D", "E"),
			rx.Just("E", "F"),
			rx.Just("F", "G"),
			rx.Pipe1(rx.Range(1, 4), DelaySubscription[int](1)),
			mapping,
		),
		"[A B C D E F 1]", "[B C D E F G 2]", ErrComplete,
	).Case(
		rx.ZipWithBuffering7(
			rx.Just("A", "B", "C"),
			rx.Just("B", "C", "D"),
			rx.Just("C", "D", "E"),
			rx.Just("D", "E", "F"),
			rx.Just("E", "F", "G"),
			rx.Just("F", "G", "H"),
			rx.Pipe1(rx.Range(1, 4), DelaySubscription[int](1)),
			mapping,
		),
		"[A B C D E F 1]", "[B C D E F G 2]", "[C D E F G H 3]", ErrComplete,
	).Case(
		rx.ZipWithBuffering7(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			rx.Just("D", "E", "F", "G"),
			rx.Just("E", "F", "G", "H"),
			rx.Just("F", "G", "H", "I"),
			rx.Pipe1(rx.Range(1, 4), DelaySubscription[int](1)),
			mapping,
		),
		"[A B C D E F 1]", "[B C D E F G 2]", "[C D E F G H 3]", ErrComplete,
	).Case(
		rx.ZipWithBuffering7(
			rx.Just("A", "B"),
			rx.Just("B", "C"),
			rx.Just("C", "D"),
			rx.Just("D", "E"),
			rx.Just("E", "F"),
			rx.Just("F", "G"),
			rx.Pipe1(
				rx.Concat(
					rx.Pipe1(rx.Range(1, 4), DelaySubscription[int](1)),
					rx.Throw[int](ErrTest),
				),
				DelaySubscription[int](1),
			),
			mapping,
		),
		"[A B C D E F 1]", "[B C D E F G 2]", ErrComplete,
	).Case(
		rx.ZipWithBuffering7(
			rx.Just("A", "B", "C"),
			rx.Just("B", "C", "D"),
			rx.Just("C", "D", "E"),
			rx.Just("D", "E", "F"),
			rx.Just("E", "F", "G"),
			rx.Just("F", "G", "H"),
			rx.Pipe1(
				rx.Concat(
					rx.Pipe1(rx.Range(1, 4), DelaySubscription[int](1)),
					rx.Throw[int](ErrTest),
				),
				DelaySubscription[int](1),
			),
			mapping,
		),
		"[A B C D E F 1]", "[B C D E F G 2]", "[C D E F G H 3]", ErrComplete,
	).Case(
		rx.ZipWithBuffering7(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			rx.Just("D", "E", "F", "G"),
			rx.Just("E", "F", "G", "H"),
			rx.Just("F", "G", "H", "I"),
			rx.Pipe1(
				rx.Concat(
					rx.Pipe1(rx.Range(1, 4), DelaySubscription[int](1)),
					rx.Throw[int](ErrTest),
				),
				DelaySubscription[int](1),
			),
			mapping,
		),
		"[A B C D E F 1]", "[B C D E F G 2]", "[C D E F G H 3]", ErrTest,
	).Case(
		rx.Pipe1(
			rx.ZipWithBuffering7(
				rx.Just("A", "B", "C", "D"),
				rx.Just("B", "C", "D", "E"),
				rx.Just("C", "D", "E", "F"),
				rx.Just("D", "E", "F", "G"),
				rx.Just("E", "F", "G", "H"),
				rx.Just("F", "G", "H", "I"),
				rx.Range(1, 5),
				mapping,
			),
			rx.OnNext(func(string) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}
