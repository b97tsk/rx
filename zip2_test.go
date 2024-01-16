package rx_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestZip2(t *testing.T) {
	t.Parallel()

	testZip2(t, rx.ConcatWith(
		func(_ context.Context, sink rx.Observer[string]) {
			sink(rx.Notification[string]{}) // For coverage.
			sink.Complete()
		},
	), ErrComplete)
	testZip2(t, rx.ConcatWith(rx.Throw[string](ErrTest)), ErrTest)
}

func testZip2(t *testing.T, op rx.Operator[string, string], err error) {
	toString := func(v1, v2 string) string {
		return fmt.Sprintf("[%v %v]", v1, v2)
	}

	NewTestSuite[string](t).Case(
		rx.Zip2(
			rx.Pipe1(rx.Just("A", "B", "C"), op),
			rx.Just("B", "C", "D", "E"),
			toString,
		),
		"[A B]", "[B C]", "[C D]", err,
	).Case(
		rx.Zip2(
			rx.Just("A", "B", "C", "D"),
			rx.Pipe1(rx.Just("B", "C", "D"), op),
			toString,
		),
		"[A B]", "[B C]", "[C D]", err,
	)
}
