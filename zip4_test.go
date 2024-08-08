package rx_test

import (
	"fmt"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestZip4(t *testing.T) {
	t.Parallel()

	testZip4(t, rx.ConcatWith(rx.Empty[string]()), ErrComplete)
	testZip4(t, rx.ConcatWith(rx.Throw[string](ErrTest)), ErrTest)
	testZip4(t, rx.ConcatWith(
		func(_ rx.Context, o rx.Observer[string]) {
			o.Emit(rx.Notification[string]{}) // For coverage.
			o.Stop(ErrTest)
		},
	), ErrTest)
}

func testZip4(t *testing.T, op rx.Operator[string, string], err error) {
	mapping := func(v1, v2, v3, v4 string) string {
		return fmt.Sprintf("[%v %v %v %v]", v1, v2, v3, v4)
	}

	NewTestSuite[string](t).Case(
		rx.Zip4(
			rx.Pipe1(rx.Just("A", "B", "C"), op),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			rx.Just("D", "E", "F", "G"),
			mapping,
		),
		"[A B C D]", "[B C D E]", "[C D E F]", err,
	).Case(
		rx.Zip4(
			rx.Just("A", "B", "C", "D"),
			rx.Pipe1(rx.Just("B", "C", "D"), op),
			rx.Just("C", "D", "E", "F"),
			rx.Just("D", "E", "F", "G"),
			mapping,
		),
		"[A B C D]", "[B C D E]", "[C D E F]", err,
	).Case(
		rx.Zip4(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Pipe1(rx.Just("C", "D", "E"), op),
			rx.Just("D", "E", "F", "G"),
			mapping,
		),
		"[A B C D]", "[B C D E]", "[C D E F]", err,
	).Case(
		rx.Zip4(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			rx.Pipe1(rx.Just("D", "E", "F"), op),
			mapping,
		),
		"[A B C D]", "[B C D E]", "[C D E F]", err,
	).Case(
		rx.Pipe1(
			rx.Zip4(
				rx.Just("A", "B", "C", "D"),
				rx.Just("B", "C", "D", "E"),
				rx.Just("C", "D", "E", "F"),
				rx.Just("D", "E", "F", "G"),
				mapping,
			),
			rx.DoOnNext(func(string) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}
