package rx_test

import (
	"fmt"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestZip3(t *testing.T) {
	t.Parallel()

	testZip3(t, rx.ConcatWith(
		func(_ rx.Context, sink rx.Observer[string]) {
			sink(rx.Notification[string]{}) // For coverage.
			sink.Complete()
		},
	), ErrComplete)
	testZip3(t, rx.ConcatWith(rx.Throw[string](ErrTest)), ErrTest)
}

func testZip3(t *testing.T, op rx.Operator[string, string], err error) {
	proj := func(v1, v2, v3 string) string {
		return fmt.Sprintf("[%v %v %v]", v1, v2, v3)
	}

	NewTestSuite[string](t).Case(
		rx.Zip3(
			rx.Pipe1(rx.Just("A", "B", "C"), op),
			rx.Just("B", "C", "D", "E"),
			rx.Just("C", "D", "E", "F"),
			proj,
		),
		"[A B C]", "[B C D]", "[C D E]", err,
	).Case(
		rx.Zip3(
			rx.Just("A", "B", "C", "D"),
			rx.Pipe1(rx.Just("B", "C", "D"), op),
			rx.Just("C", "D", "E", "F"),
			proj,
		),
		"[A B C]", "[B C D]", "[C D E]", err,
	).Case(
		rx.Zip3(
			rx.Just("A", "B", "C", "D"),
			rx.Just("B", "C", "D", "E"),
			rx.Pipe1(rx.Just("C", "D", "E"), op),
			proj,
		),
		"[A B C]", "[B C D]", "[C D E]", err,
	).Case(
		rx.Pipe1(
			rx.Zip3(
				rx.Just("A", "B", "C", "D"),
				rx.Just("B", "C", "D", "E"),
				rx.Just("C", "D", "E", "F"),
				proj,
			),
			rx.OnNext(func(string) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}
