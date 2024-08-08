package rx_test

import (
	"fmt"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestZip5(t *testing.T) {
	t.Parallel()

	testZip5(t, rx.ConcatWith(rx.Empty[string]()), ErrComplete)
	testZip5(t, rx.ConcatWith(rx.Throw[string](ErrTest)), ErrTest)
	testZip5(t, rx.ConcatWith(
		func(_ rx.Context, o rx.Observer[string]) {
			o.Emit(rx.Notification[string]{}) // For coverage.
			o.Stop(ErrTest)
		},
	), ErrTest)
}

func testZip5(t *testing.T, op rx.Operator[string, string], err error) {
	mapping := func(v1, v2, v3, v4, v5 string) string {
		return fmt.Sprintf("[%v %v %v %v %v]", v1, v2, v3, v4, v5)
	}

	NewTestSuite[string](t).Case(
		rx.Zip5(
			rx.Pipe1(rx.Just("A", "B", "C"), op),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			rx.Just("D", "E", "F", "G"),
			rx.Just("E", "F", "G", "H"),
			mapping,
		),
		"[A B C D E]", "[B C D E F]", "[C D E F G]", err,
	).Case(
		rx.Zip5(
			rx.Just("A", "B", "C", "D"),
			rx.Pipe1(rx.Just("B", "C", "D"), op),
			rx.Just("C", "D", "E", "F"),
			rx.Just("D", "E", "F", "G"),
			rx.Just("E", "F", "G", "H"),
			mapping,
		),
		"[A B C D E]", "[B C D E F]", "[C D E F G]", err,
	).Case(
		rx.Zip5(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Pipe1(rx.Just("C", "D", "E"), op),
			rx.Just("D", "E", "F", "G"),
			rx.Just("E", "F", "G", "H"),
			mapping,
		),
		"[A B C D E]", "[B C D E F]", "[C D E F G]", err,
	).Case(
		rx.Zip5(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			rx.Pipe1(rx.Just("D", "E", "F"), op),
			rx.Just("E", "F", "G", "H"),
			mapping,
		),
		"[A B C D E]", "[B C D E F]", "[C D E F G]", err,
	).Case(
		rx.Zip5(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			rx.Just("D", "E", "F", "G"),
			rx.Pipe1(rx.Just("E", "F", "G"), op),
			mapping,
		),
		"[A B C D E]", "[B C D E F]", "[C D E F G]", err,
	).Case(
		rx.Pipe1(
			rx.Zip5(
				rx.Just("A", "B", "C", "D"),
				rx.Just("B", "C", "D", "E"),
				rx.Just("C", "D", "E", "F"),
				rx.Just("D", "E", "F", "G"),
				rx.Just("E", "F", "G", "H"),
				mapping,
			),
			rx.DoOnNext(func(string) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}
