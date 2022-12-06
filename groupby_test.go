package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestGroupBy(t *testing.T) {
	t.Parallel()

	count := rx.Compose2(
		rx.GroupBy(
			func(v string) string { return v },
			func() rx.Subject[string] { return rx.MulticastReplay[string](nil) },
		),
		rx.ConcatMap(
			func(group rx.Pair[string, rx.Observable[string]]) rx.Observable[rx.Pair[string, int]] {
				return rx.Pipe2(
					group.Value,
					rx.Reduce(0, func(n int, _ string) int { return n + 1 }),
					rx.Map(
						func(v int) rx.Pair[string, int] {
							return rx.NewPair(group.Key, v)
						},
					),
				)
			},
		),
	)

	NewTestSuite[string](t).Case(
		rx.Pipe3(
			rx.Just("A", "B", "B", "A", "C", "C", "D", "A"),
			AddLatencyToValues[string](0, 1),
			count,
			ToString[rx.Pair[string, int]](),
		),
		"{A 3}", "{B 2}", "{C 2}", "{D 1}", ErrCompleted,
	).Case(
		rx.Pipe2(
			rx.Throw[string](ErrTest),
			count,
			ToString[rx.Pair[string, int]](),
		),
		ErrTest,
	)
}
