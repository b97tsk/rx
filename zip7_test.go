package rx_test

import (
	"fmt"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestZip7(t *testing.T) {
	t.Parallel()

	testZip7(t, rx.ConcatWith(
		func(_ rx.Context, sink rx.Observer[string]) {
			sink(rx.Notification[string]{}) // For coverage.
			sink.Complete()
		},
	), ErrComplete)
	testZip7(t, rx.ConcatWith(rx.Throw[string](ErrTest)), ErrTest)
}

func testZip7(t *testing.T, op rx.Operator[string, string], err error) {
	mapping := func(v1, v2, v3, v4, v5, v6, v7 string) string {
		return fmt.Sprintf("[%v %v %v %v %v %v %v]", v1, v2, v3, v4, v5, v6, v7)
	}

	NewTestSuite[string](t).Case(
		rx.Zip7(
			rx.Pipe1(rx.Just("A", "B", "C"), op),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			rx.Just("D", "E", "F", "G"),
			rx.Just("E", "F", "G", "H"),
			rx.Just("F", "G", "H", "I"),
			rx.Just("G", "H", "I", "J"),
			mapping,
		),
		"[A B C D E F G]", "[B C D E F G H]", "[C D E F G H I]", err,
	).Case(
		rx.Zip7(
			rx.Just("A", "B", "C", "D"),
			rx.Pipe1(rx.Just("B", "C", "D"), op),
			rx.Just("C", "D", "E", "F"),
			rx.Just("D", "E", "F", "G"),
			rx.Just("E", "F", "G", "H"),
			rx.Just("F", "G", "H", "I"),
			rx.Just("G", "H", "I", "J"),
			mapping,
		),
		"[A B C D E F G]", "[B C D E F G H]", "[C D E F G H I]", err,
	).Case(
		rx.Zip7(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Pipe1(rx.Just("C", "D", "E"), op),
			rx.Just("D", "E", "F", "G"),
			rx.Just("E", "F", "G", "H"),
			rx.Just("F", "G", "H", "I"),
			rx.Just("G", "H", "I", "J"),
			mapping,
		),
		"[A B C D E F G]", "[B C D E F G H]", "[C D E F G H I]", err,
	).Case(
		rx.Zip7(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			rx.Pipe1(rx.Just("D", "E", "F"), op),
			rx.Just("E", "F", "G", "H"),
			rx.Just("F", "G", "H", "I"),
			rx.Just("G", "H", "I", "J"),
			mapping,
		),
		"[A B C D E F G]", "[B C D E F G H]", "[C D E F G H I]", err,
	).Case(
		rx.Zip7(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			rx.Just("D", "E", "F", "G"),
			rx.Pipe1(rx.Just("E", "F", "G"), op),
			rx.Just("F", "G", "H", "I"),
			rx.Just("G", "H", "I", "J"),
			mapping,
		),
		"[A B C D E F G]", "[B C D E F G H]", "[C D E F G H I]", err,
	).Case(
		rx.Zip7(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			rx.Just("D", "E", "F", "G"),
			rx.Just("E", "F", "G", "H"),
			rx.Pipe1(rx.Just("F", "G", "H"), op),
			rx.Just("G", "H", "I", "J"),
			mapping,
		),
		"[A B C D E F G]", "[B C D E F G H]", "[C D E F G H I]", err,
	).Case(
		rx.Zip7(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			rx.Just("D", "E", "F", "G"),
			rx.Just("E", "F", "G", "H"),
			rx.Just("F", "G", "H", "I"),
			rx.Pipe1(rx.Just("G", "H", "I"), op),
			mapping,
		),
		"[A B C D E F G]", "[B C D E F G H]", "[C D E F G H I]", err,
	).Case(
		rx.Pipe1(
			rx.Zip7(
				rx.Just("A", "B", "C", "D"),
				rx.Just("B", "C", "D", "E"),
				rx.Just("C", "D", "E", "F"),
				rx.Just("D", "E", "F", "G"),
				rx.Just("E", "F", "G", "H"),
				rx.Just("F", "G", "H", "I"),
				rx.Just("G", "H", "I", "J"),
				mapping,
			),
			rx.OnNext(func(string) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}
