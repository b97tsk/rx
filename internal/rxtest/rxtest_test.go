package rxtest_test

import (
	"context"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestSuccess(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).WithContext(context.Background()).Case(
		rx.Pipe4(
			rx.Just(42),
			AddLatencyToValues[int](0, 1),
			AddLatencyToNotifications[int](0, 1),
			DelaySubscription[int](1),
			ToString[int](),
		),
		"42",
		ErrCompleted,
	).Case(
		rx.Throw[string](ErrTest),
		ErrTest,
	)
}

func TestFailure(t *testing.T) {
	t.Parallel()

	failtest(t, rx.Just(ErrTest))
	failtest(t, rx.Throw[string](ErrTest))
	failtest(t, rx.Just(ErrTest), ErrCompleted, ErrTest)
	failtest(t, rx.Throw[string](ErrTest), ErrCompleted, ErrTest)
}

func failtest[T any](tb testing.TB, obs rx.Observable[T], output ...any) {
	tb.Helper()

	fs := &failsafe{TB: tb}

	NewTestSuite[T](fs).Case(obs, output...)

	if !fs.failed {
		tb.FailNow()
	}
}

type failsafe struct {
	testing.TB
	failed bool
}

func (fs *failsafe) Fail() {
	fs.failed = true
}
