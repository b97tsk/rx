package rx_test

import (
	"fmt"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestZip9(t *testing.T) {
	t.Parallel()

	testZip9(t, rx.ConcatWith(
		func(_ rx.Context, o rx.Observer[string]) {
			o.Emit(rx.Notification[string]{}) // For coverage.
			o.Complete()
		},
	), ErrComplete)
	testZip9(t, rx.ConcatWith(rx.Throw[string](ErrTest)), ErrTest)
}

func testZip9(t *testing.T, op rx.Operator[string, string], err error) {
	mapping := func(v1, v2, v3, v4, v5, v6, v7, v8, v9 string) string {
		return fmt.Sprintf("[%v %v %v %v %v %v %v %v %v]", v1, v2, v3, v4, v5, v6, v7, v8, v9)
	}

	NewTestSuite[string](t).Case(
		rx.Zip9(
			rx.Pipe1(rx.Just("A", "B", "C"), op),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			rx.Just("D", "E", "F", "G"),
			rx.Just("E", "F", "G", "H"),
			rx.Just("F", "G", "H", "I"),
			rx.Just("G", "H", "I", "J"),
			rx.Just("H", "I", "J", "K"),
			rx.Just("I", "J", "K", "L"),
			mapping,
		),
		"[A B C D E F G H I]", "[B C D E F G H I J]", "[C D E F G H I J K]", err,
	).Case(
		rx.Zip9(
			rx.Just("A", "B", "C", "D"),
			rx.Pipe1(rx.Just("B", "C", "D"), op),
			rx.Just("C", "D", "E", "F"),
			rx.Just("D", "E", "F", "G"),
			rx.Just("E", "F", "G", "H"),
			rx.Just("F", "G", "H", "I"),
			rx.Just("G", "H", "I", "J"),
			rx.Just("H", "I", "J", "K"),
			rx.Just("I", "J", "K", "L"),
			mapping,
		),
		"[A B C D E F G H I]", "[B C D E F G H I J]", "[C D E F G H I J K]", err,
	).Case(
		rx.Zip9(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Pipe1(rx.Just("C", "D", "E"), op),
			rx.Just("D", "E", "F", "G"),
			rx.Just("E", "F", "G", "H"),
			rx.Just("F", "G", "H", "I"),
			rx.Just("G", "H", "I", "J"),
			rx.Just("H", "I", "J", "K"),
			rx.Just("I", "J", "K", "L"),
			mapping,
		),
		"[A B C D E F G H I]", "[B C D E F G H I J]", "[C D E F G H I J K]", err,
	).Case(
		rx.Zip9(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			rx.Pipe1(rx.Just("D", "E", "F"), op),
			rx.Just("E", "F", "G", "H"),
			rx.Just("F", "G", "H", "I"),
			rx.Just("G", "H", "I", "J"),
			rx.Just("H", "I", "J", "K"),
			rx.Just("I", "J", "K", "L"),
			mapping,
		),
		"[A B C D E F G H I]", "[B C D E F G H I J]", "[C D E F G H I J K]", err,
	).Case(
		rx.Zip9(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			rx.Just("D", "E", "F", "G"),
			rx.Pipe1(rx.Just("E", "F", "G"), op),
			rx.Just("F", "G", "H", "I"),
			rx.Just("G", "H", "I", "J"),
			rx.Just("H", "I", "J", "K"),
			rx.Just("I", "J", "K", "L"),
			mapping,
		),
		"[A B C D E F G H I]", "[B C D E F G H I J]", "[C D E F G H I J K]", err,
	).Case(
		rx.Zip9(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			rx.Just("D", "E", "F", "G"),
			rx.Just("E", "F", "G", "H"),
			rx.Pipe1(rx.Just("F", "G", "H"), op),
			rx.Just("G", "H", "I", "J"),
			rx.Just("H", "I", "J", "K"),
			rx.Just("I", "J", "K", "L"),
			mapping,
		),
		"[A B C D E F G H I]", "[B C D E F G H I J]", "[C D E F G H I J K]", err,
	).Case(
		rx.Zip9(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			rx.Just("D", "E", "F", "G"),
			rx.Just("E", "F", "G", "H"),
			rx.Just("F", "G", "H", "I"),
			rx.Pipe1(rx.Just("G", "H", "I"), op),
			rx.Just("H", "I", "J", "K"),
			rx.Just("I", "J", "K", "L"),
			mapping,
		),
		"[A B C D E F G H I]", "[B C D E F G H I J]", "[C D E F G H I J K]", err,
	).Case(
		rx.Zip9(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			rx.Just("D", "E", "F", "G"),
			rx.Just("E", "F", "G", "H"),
			rx.Just("F", "G", "H", "I"),
			rx.Just("G", "H", "I", "J"),
			rx.Pipe1(rx.Just("H", "I", "J"), op),
			rx.Just("I", "J", "K", "L"),
			mapping,
		),
		"[A B C D E F G H I]", "[B C D E F G H I J]", "[C D E F G H I J K]", err,
	).Case(
		rx.Zip9(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			rx.Just("D", "E", "F", "G"),
			rx.Just("E", "F", "G", "H"),
			rx.Just("F", "G", "H", "I"),
			rx.Just("G", "H", "I", "J"),
			rx.Just("H", "I", "J", "K"),
			rx.Pipe1(rx.Just("I", "J", "K"), op),
			mapping,
		),
		"[A B C D E F G H I]", "[B C D E F G H I J]", "[C D E F G H I J K]", err,
	).Case(
		rx.Pipe1(
			rx.Zip9(
				rx.Just("A", "B", "C", "D"),
				rx.Just("B", "C", "D", "E"),
				rx.Just("C", "D", "E", "F"),
				rx.Just("D", "E", "F", "G"),
				rx.Just("E", "F", "G", "H"),
				rx.Just("F", "G", "H", "I"),
				rx.Just("G", "H", "I", "J"),
				rx.Just("H", "I", "J", "K"),
				rx.Just("I", "J", "K", "L"),
				mapping,
			),
			rx.DoOnNext(func(string) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}
