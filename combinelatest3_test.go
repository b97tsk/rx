package rx_test

import (
	"fmt"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestCombineLatest3(t *testing.T) {
	t.Parallel()

	toString := func(v1, v2, v3 string) string {
		return fmt.Sprintf("[%v %v %v]", v1, v2, v3)
	}

	NewTestSuite[string](t).Case(
		rx.CombineLatest3(
			rx.Pipe(rx.Just("A", "D"), AddLatencyToValues[string](1, 3)),
			rx.Pipe(rx.Just("B", "E"), AddLatencyToValues[string](2, 3)),
			rx.Pipe(rx.Just("C", "F"), AddLatencyToValues[string](3, 3)),
			toString,
		),
		"[A B C]", "[D B C]", "[D E C]", "[D E F]", ErrCompleted,
	).Case(
		rx.CombineLatest3(
			rx.Throw[string](ErrTest),
			rx.Throw[string](ErrTest),
			rx.Throw[string](ErrTest),
			toString,
		),
		ErrTest,
	)
}
