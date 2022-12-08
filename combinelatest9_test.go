package rx_test

import (
	"fmt"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestCombineLatest9(t *testing.T) {
	t.Parallel()

	toString := func(v1, v2, v3, v4, v5, v6, v7, v8, v9 string) string {
		return fmt.Sprintf("[%v %v %v %v %v %v %v %v %v]", v1, v2, v3, v4, v5, v6, v7, v8, v9)
	}

	NewTestSuite[string](t).Case(
		rx.CombineLatest9(
			rx.Pipe(rx.Just("A", "J"), AddLatencyToValues[string](1, 9)),
			rx.Pipe(rx.Just("B", "K"), AddLatencyToValues[string](2, 9)),
			rx.Pipe(rx.Just("C", "L"), AddLatencyToValues[string](3, 9)),
			rx.Pipe(rx.Just("D", "M"), AddLatencyToValues[string](4, 9)),
			rx.Pipe(rx.Just("E", "N"), AddLatencyToValues[string](5, 9)),
			rx.Pipe(rx.Just("F", "O"), AddLatencyToValues[string](6, 9)),
			rx.Pipe(rx.Just("G", "P"), AddLatencyToValues[string](7, 9)),
			rx.Pipe(rx.Just("H", "Q"), AddLatencyToValues[string](8, 9)),
			rx.Pipe(rx.Just("I", "R"), AddLatencyToValues[string](9, 9)),
			toString,
		),
		"[A B C D E F G H I]", "[J B C D E F G H I]", "[J K C D E F G H I]",
		"[J K L D E F G H I]", "[J K L M E F G H I]", "[J K L M N F G H I]",
		"[J K L M N O G H I]", "[J K L M N O P H I]", "[J K L M N O P Q I]",
		"[J K L M N O P Q R]",
		ErrCompleted,
	).Case(
		rx.CombineLatest9(
			rx.Throw[string](ErrTest),
			rx.Throw[string](ErrTest),
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
