package rx_test

import (
	"context"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestOnErrorResumeNext(t *testing.T) {
	NewTestSuite(t).Case(
		rx.OnErrorResumeNext(
			rx.Range(1, 4),
			rx.Range(4, 7),
			rx.Range(7, 10),
		),
		1, 2, 3, 4, 5, 6, 7, 8, 9, Completed,
	).Case(
		rx.OnErrorResumeNext(
			rx.Throw(ErrTest),
			rx.Range(4, 7),
			rx.Range(7, 10),
		),
		4, 5, 6, 7, 8, 9, Completed,
	).Case(
		rx.OnErrorResumeNext(
			rx.Range(1, 4),
			rx.Throw(ErrTest),
			rx.Range(7, 10),
		),
		1, 2, 3, 7, 8, 9, Completed,
	).Case(
		rx.OnErrorResumeNext(
			rx.Range(1, 4),
			rx.Range(4, 7),
			rx.Throw(ErrTest),
		),
		1, 2, 3, 4, 5, 6, Completed,
	).Case(
		rx.OnErrorResumeNext(),
		Completed,
	).TestAll()

	ctx, cancel := context.WithTimeout(context.Background(), Step(1))
	defer cancel()

	rx.OnErrorResumeNext(
		func(_ context.Context, sink rx.Observer) {
			_ = rx.Timer(Step(2)).BlockingSubscribe(context.Background(), sink)
		},
		func(context.Context, rx.Observer) {
			t.Fatal("should not happen")
		},
	).Subscribe(ctx, rx.Noop)
}
