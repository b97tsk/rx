package rx_test

import (
	"fmt"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestCombineLatest7(t *testing.T) {
	t.Parallel()

	toString := func(v1, v2, v3, v4, v5, v6, v7 string) string {
		return fmt.Sprintf("[%v %v %v %v %v %v %v]", v1, v2, v3, v4, v5, v6, v7)
	}

	NewTestSuite[string](t).Case(
		rx.CombineLatest7(
			rx.Pipe(rx.Just("A", "H"), AddLatencyToValues[string](1, 7)),
			rx.Pipe(rx.Just("B", "I"), AddLatencyToValues[string](2, 7)),
			rx.Pipe(rx.Just("C", "J"), AddLatencyToValues[string](3, 7)),
			rx.Pipe(rx.Just("D", "K"), AddLatencyToValues[string](4, 7)),
			rx.Pipe(rx.Just("E", "L"), AddLatencyToValues[string](5, 7)),
			rx.Pipe(rx.Just("F", "M"), AddLatencyToValues[string](6, 7)),
			rx.Pipe(rx.Just("G", "N"), AddLatencyToValues[string](7, 7)),
			toString,
		),
		"[A B C D E F G]", "[H B C D E F G]", "[H I C D E F G]", "[H I J D E F G]",
		"[H I J K E F G]", "[H I J K L F G]", "[H I J K L M G]", "[H I J K L M N]",
		ErrCompleted,
	).Case(
		rx.CombineLatest7(
			rx.Throw[string](ErrTest),
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
