package rxtest_test

import (
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestSuccess(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).WithContext(rx.NewBackgroundContext()).Case(
		rx.Pipe4(
			rx.Just(42),
			AddLatencyToValues[int](0, 1),
			AddLatencyToNotifications[int](0, 1),
			DelaySubscription[int](1),
			ToString[int](),
		),
		"42",
		ErrComplete,
	).Case(
		rx.Throw[string](ErrTest),
		ErrTest,
	).Case(
		rx.Oops[string](ErrTest),
		rx.ErrOops, ErrTest,
	)
}

func TestFailure(t *testing.T) {
	t.Parallel()

	failtest(t, rx.Just(42))
	failtest(t, rx.Just(42), ErrTest, ErrTest, ErrTest)
	failtest(t, rx.Throw[string](ErrTest))
	failtest(t, rx.Throw[string](ErrTest), ErrComplete, ErrComplete)
	failtest(t, rx.Oops[string](ErrTest))
	failtest(t, rx.Oops[string](ErrTest), ErrComplete, ErrComplete, ErrComplete)
	failtest(t, func(c rx.Context, sink rx.Observer[string]) {
		c.Go(func() { time.Sleep(8 * time.Second) })
		sink.Complete()
	}, ErrComplete)
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

func (fs *failsafe) Error(args ...any) {
	fs.failed = true
	fs.Log(args...)
}

func (fs *failsafe) Errorf(format string, args ...any) {
	fs.failed = true
	fs.Logf(format, args...)
}

func (fs *failsafe) Fail() {
	fs.failed = true
}
