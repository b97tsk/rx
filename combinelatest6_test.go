package rx_test

import (
	"fmt"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestCombineLatest6(t *testing.T) {
	t.Parallel()

	toString := func(v1, v2, v3, v4, v5, v6 string) string {
		return fmt.Sprintf("[%v %v %v %v %v %v]", v1, v2, v3, v4, v5, v6)
	}

	NewTestSuite[string](t).Case(
		rx.CombineLatest6(
			rx.Pipe(rx.Just("A", "G"), AddLatencyToValues[string](1, 6)),
			rx.Pipe(rx.Just("B", "H"), AddLatencyToValues[string](2, 6)),
			rx.Pipe(rx.Just("C", "I"), AddLatencyToValues[string](3, 6)),
			rx.Pipe(rx.Just("D", "J"), AddLatencyToValues[string](4, 6)),
			rx.Pipe(rx.Just("E", "K"), AddLatencyToValues[string](5, 6)),
			rx.Pipe(rx.Just("F", "L"), AddLatencyToValues[string](6, 6)),
			toString,
		),
		"[A B C D E F]", "[G B C D E F]", "[G H C D E F]", "[G H I D E F]",
		"[G H I J E F]", "[G H I J K F]", "[G H I J K L]",
		ErrCompleted,
	).Case(
		rx.CombineLatest6(
			rx.Throw[string](ErrTest),
			rx.Throw[string](ErrTest),
			rx.Throw[string](ErrTest),
			rx.Throw[string](ErrTest),
			rx.Throw[string](ErrTest),
			rx.Throw[string](ErrTest),
			toString,
		),
		ErrTest,
	)
}
