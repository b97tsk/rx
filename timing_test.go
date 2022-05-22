package rx_test

import (
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestTicker(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe2(
			rx.Ticker(Step(1)),
			rx.Scan(-1, func(i int, _ time.Time) int { return i + 1 }),
			rx.Take[int](3),
		),
		0, 1, 2, ErrCompleted,
	)
}

func TestTimer(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe(
			rx.Timer(Step(1)),
			rx.MapTo[time.Time](42),
		),
		42, ErrCompleted,
	)
}
