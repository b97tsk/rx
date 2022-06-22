package rx_test

import (
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestConnect(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe2(
			rx.Ticker(Step(1)),
			rx.Scan(-1, func(i int, _ time.Time) int { return i + 1 }),
			rx.Connect(
				func(source rx.Observable[int]) rx.Observable[string] {
					return rx.Pipe2(
						rx.Zip2(
							rx.Pipe(source, rx.Take[int](4)),
							rx.Pipe2(source, rx.Skip[int](4), rx.Take[int](4)),
							func(v1, v2 int) int { return v1 * v2 },
						),
						rx.ToSlice[int](),
						ToString[[]int](),
					)
				},
			).WithConnector(rx.Multicast[int]).AsOperator(),
		),
		"[0 5 12 21]", ErrCompleted,
	)
}
