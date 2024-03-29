package rx_test

import (
	"fmt"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestZip6(t *testing.T) {
	t.Parallel()

	testZip6(t, rx.ConcatWith(
		func(_ rx.Context, sink rx.Observer[string]) {
			sink(rx.Notification[string]{}) // For coverage.
			sink.Complete()
		},
	), ErrComplete)
	testZip6(t, rx.ConcatWith(rx.Throw[string](ErrTest)), ErrTest)
}

func testZip6(t *testing.T, op rx.Operator[string, string], err error) {
	mapping := func(v1, v2, v3, v4, v5, v6 string) string {
		return fmt.Sprintf("[%v %v %v %v %v %v]", v1, v2, v3, v4, v5, v6)
	}

	NewTestSuite[string](t).Case(
		rx.Zip6(
			rx.Pipe1(rx.Just("A", "B", "C"), op),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			rx.Just("D", "E", "F", "G"),
			rx.Just("E", "F", "G", "H"),
			rx.Just("F", "G", "H", "I"),
			mapping,
		),
		"[A B C D E F]", "[B C D E F G]", "[C D E F G H]", err,
	).Case(
		rx.Zip6(
			rx.Just("A", "B", "C", "D"),
			rx.Pipe1(rx.Just("B", "C", "D"), op),
			rx.Just("C", "D", "E", "F"),
			rx.Just("D", "E", "F", "G"),
			rx.Just("E", "F", "G", "H"),
			rx.Just("F", "G", "H", "I"),
			mapping,
		),
		"[A B C D E F]", "[B C D E F G]", "[C D E F G H]", err,
	).Case(
		rx.Zip6(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Pipe1(rx.Just("C", "D", "E"), op),
			rx.Just("D", "E", "F", "G"),
			rx.Just("E", "F", "G", "H"),
			rx.Just("F", "G", "H", "I"),
			mapping,
		),
		"[A B C D E F]", "[B C D E F G]", "[C D E F G H]", err,
	).Case(
		rx.Zip6(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			rx.Pipe1(rx.Just("D", "E", "F"), op),
			rx.Just("E", "F", "G", "H"),
			rx.Just("F", "G", "H", "I"),
			mapping,
		),
		"[A B C D E F]", "[B C D E F G]", "[C D E F G H]", err,
	).Case(
		rx.Zip6(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			rx.Just("D", "E", "F", "G"),
			rx.Pipe1(rx.Just("E", "F", "G"), op),
			rx.Just("F", "G", "H", "I"),
			mapping,
		),
		"[A B C D E F]", "[B C D E F G]", "[C D E F G H]", err,
	).Case(
		rx.Zip6(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			rx.Just("D", "E", "F", "G"),
			rx.Just("E", "F", "G", "H"),
			rx.Pipe1(rx.Just("F", "G", "H"), op),
			mapping,
		),
		"[A B C D E F]", "[B C D E F G]", "[C D E F G H]", err,
	).Case(
		rx.Pipe1(
			rx.Zip6(
				rx.Just("A", "B", "C", "D"),
				rx.Just("B", "C", "D", "E"),
				rx.Just("C", "D", "E", "F"),
				rx.Just("D", "E", "F", "G"),
				rx.Just("E", "F", "G", "H"),
				rx.Just("F", "G", "H", "I"),
				mapping,
			),
			rx.DoOnNext(func(string) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}
