package rx_test

import (
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func congestOnNext(v int) {
	n := 4

	if v == 1 {
		// The first value arrives at odd step.
		// Decrease n to make sure follow-ups called at even steps so as to
		// pass the tests.
		n--
	}

	time.Sleep(Step(n))
}

func TestCongestBlock(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe3(
			rx.Range(1, 8),
			AddLatencyToValues[int](1, 2),
			rx.CongestBlock[int](3),
			rx.DoOnNext(congestOnNext),
		),
		1, 2, 3, 4, 5, 6, 7, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Empty[int](),
			rx.CongestBlock[int](3),
		),
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Throw[int](ErrTest),
			rx.CongestBlock[int](3),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.Empty[int](),
			rx.CongestBlock[int](0),
		),
		rx.ErrOops, "CongestBlock: capacity < 1",
	)
}

func TestCongestDropLatest(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe3(
			rx.Range(1, 8),
			AddLatencyToValues[int](1, 2),
			rx.CongestDropLatest[int](3),
			rx.DoOnNext(congestOnNext),
		),
		1, 2, 3, 4, 5, 6, 7, ErrComplete,
	).Case(
		rx.Pipe3(
			rx.Range(1, 11),
			AddLatencyToValues[int](1, 2),
			rx.CongestDropLatest[int](3),
			rx.DoOnNext(congestOnNext),
		),
		1, 2, 3, 4, 5, 6, 7, 9, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Empty[int](),
			rx.CongestDropLatest[int](3),
		),
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Throw[int](ErrTest),
			rx.CongestDropLatest[int](3),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.Empty[int](),
			rx.CongestDropLatest[int](0),
		),
		rx.ErrOops, "CongestDropLatest: capacity < 1",
	)
}

func TestCongestDropOldest(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe3(
			rx.Range(1, 8),
			AddLatencyToValues[int](1, 2),
			rx.CongestDropOldest[int](3),
			rx.DoOnNext(congestOnNext),
		),
		1, 2, 3, 4, 5, 6, 7, ErrComplete,
	).Case(
		rx.Pipe3(
			rx.Range(1, 11),
			AddLatencyToValues[int](1, 2),
			rx.CongestDropOldest[int](3),
			rx.DoOnNext(congestOnNext),
		),
		1, 2, 3, 4, 6, 8, 9, 10, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Empty[int](),
			rx.CongestDropOldest[int](3),
		),
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Throw[int](ErrTest),
			rx.CongestDropOldest[int](3),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.Empty[int](),
			rx.CongestDropOldest[int](0),
		),
		rx.ErrOops, "CongestDropOldest: capacity < 1",
	)
}

func TestCongestError(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe3(
			rx.Range(1, 8),
			AddLatencyToValues[int](1, 2),
			rx.CongestError[int](3),
			rx.DoOnNext(congestOnNext),
		),
		1, 2, 3, 4, 5, 6, 7, ErrComplete,
	).Case(
		rx.Pipe3(
			rx.Range(1, 11),
			AddLatencyToValues[int](1, 2),
			rx.CongestError[int](3),
			rx.DoOnNext(congestOnNext),
		),
		1, 2, 3, 4, rx.ErrBufferOverflow,
	).Case(
		rx.Pipe1(
			rx.Empty[int](),
			rx.CongestError[int](3),
		),
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Throw[int](ErrTest),
			rx.CongestError[int](3),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.Empty[int](),
			rx.CongestError[int](0),
		),
		rx.ErrOops, "CongestError: capacity < 1",
	)
}
