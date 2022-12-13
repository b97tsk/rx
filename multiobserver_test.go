package rx_test

import (
	"context"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestMultiObserver(t *testing.T) {
	t.Parallel()

	m := rx.Multicast[string]()

	rx.Pipe1(
		rx.Just("A", "B", "C"),
		AddLatencyToValues[string](1, 1),
	).Subscribe(context.Background(), m.Observer.ElementsOnly)

	NewTestSuite[string](t).Case(
		rx.Pipe1(
			m.Observable,
			rx.Take[string](2),
		),
		"A", "B", ErrCompleted,
	)

	time.Sleep(Step(5))

	ctx, cancel := context.WithTimeout(context.Background(), Step(1))
	defer cancel()

	rx.Pipe1(
		m.Observable,
		rx.DoOnNext(
			func(string) {
				time.Sleep(Step(2))
			},
		),
	).Subscribe(ctx, rx.Noop[string])

	m.Observer.Next("D")

	rx.Pipe1(
		m.Observable,
		rx.DoOnNext(
			func(string) {
				m.Observable.Subscribe(context.Background(), rx.Noop[string])
			},
		),
	).Subscribe(context.Background(), rx.Noop[string])

	m.Observer.Next("E")
	m.Observer.Complete()

	time.Sleep(Step(5))
}
