package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestGroupBy(t *testing.T) {
	t.Parallel()

	source := rx.Pipe1(
		rx.Just("A", "B", "B", "A", "C", "C", "D", "A"),
		AddLatencyToValues[string](0, 1),
	)

	group := rx.GroupBy(
		func(v string) string { return v },
		rx.MulticastBufferAll[string],
	)

	count := rx.ConcatMap(
		func(g rx.Pair[string, rx.Observable[string]]) rx.Observable[rx.Pair[string, int]] {
			return rx.Pipe2(
				g.Value,
				rx.Reduce(0, func(n int, _ string) int { return n + 1 }),
				rx.Map(
					func(v int) rx.Pair[string, int] {
						return rx.NewPair(g.Key, v)
					},
				),
			)
		},
	).WithBuffering()

	tostring := ToString[rx.Pair[string, int]]()

	NewTestSuite[string](t).Case(
		rx.Pipe3(source, group, count, tostring),
		"{A 3}", "{B 2}", "{C 2}", "{D 1}", ErrComplete,
	).Case(
		rx.Pipe3(
			rx.Throw[string](ErrTest),
			group, count, tostring,
		),
		ErrTest,
	).Case(
		rx.Pipe3(
			rx.Oops[string](ErrTest),
			group, count, tostring,
		),
		rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe3(
			source,
			rx.GroupBy(
				func(v string) string { panic(ErrTest) },
				rx.MulticastBufferAll[string],
			),
			count,
			tostring,
		),
		rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe3(
			source,
			rx.GroupBy(
				func(v string) string { return v },
				func() rx.Subject[string] { panic(ErrTest) },
			),
			count,
			tostring,
		),
		rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe2(
			source,
			group,
			rx.ConcatMap(
				func(g rx.Pair[string, rx.Observable[string]]) rx.Observable[string] {
					return rx.Pipe2(
						g.Value,
						rx.Discard[string](),
						rx.DoOnComplete[string](func() { panic(ErrTest) }),
					)
				},
			).WithBuffering(),
		),
		rx.ErrOops, ErrTest,
	)
}
