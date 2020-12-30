package rxtest_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestSuccess(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Just(ErrTest).Pipe(
			AddLatencyToValues(0, 1),
			AddLatencyToNotifications(0, 1),
			DelaySubscription(1),
			ToString(),
		),
		ErrTest.Error(),
		Completed,
	).Case(
		rx.Throw(ErrTest),
		ErrTest,
	).TestAll()
}

func TestFailure(t *testing.T) {
	failtest(t, rx.Just(ErrTest))
	failtest(t, rx.Throw(ErrTest))
	failtest(t, rx.Just(ErrTest), Completed, ErrTest)
	failtest(t, rx.Throw(ErrTest), Completed, ErrTest)
}

func failtest(t testing.TB, source rx.Observable, output ...interface{}) {
	fs := &failsafe{TB: t}
	Test(fs, source, output...)

	if !fs.failed {
		t.FailNow()
	}
}

type failsafe struct {
	testing.TB
	failed bool
}

func (fs *failsafe) Fail() {
	fs.failed = true
}
