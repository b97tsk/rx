package rx_test

import (
	"fmt"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestCombineLatest4(t *testing.T) {
	t.Parallel()

	toString := func(v1, v2, v3, v4 string) string {
		return fmt.Sprintf("[%v %v %v %v]", v1, v2, v3, v4)
	}

	NewTestSuite[string](t).Case(
		rx.CombineLatest4(
			rx.Pipe(rx.Just("A", "E"), AddLatencyToValues[string](1, 4)),
			rx.Pipe(rx.Just("B", "F"), AddLatencyToValues[string](2, 4)),
			rx.Pipe(rx.Just("C", "G"), AddLatencyToValues[string](3, 4)),
			rx.Pipe(rx.Just("D", "H"), AddLatencyToValues[string](4, 4)),
			toString,
		),
		"[A B C D]", "[E B C D]", "[E F C D]", "[E F G D]", "[E F G H]", ErrCompleted,
	).Case(
		rx.CombineLatest4(
			rx.Throw[string](ErrTest),
			rx.Throw[string](ErrTest),
			rx.Throw[string](ErrTest),
			rx.Throw[string](ErrTest),
			toString,
		),
		ErrTest,
	)
}
