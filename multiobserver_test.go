package rx_test

import (
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestMultiObserver(t *testing.T) {
	t.Parallel()

	m := rx.Multicast[string]()
	defer m.Complete()

	rx.Pipe1(
		rx.Just("A", "B", "C"),
		AddLatencyToValues[string](1, 1),
	).Subscribe(rx.NewBackgroundContext(), m.ElementsOnly)

	NewTestSuite[string](t).Case(
		rx.Pipe1(
			m.Observable,
			rx.Take[string](2),
		),
		"A", "B", ErrComplete,
	)

	time.Sleep(Step(5))

	ctx, cancel := rx.NewBackgroundContext().WithTimeout(Step(1))
	defer cancel()

	rx.Pipe1(
		m.Observable,
		rx.DoOnNext(
			func(string) {
				time.Sleep(Step(2))
			},
		),
	).Subscribe(ctx, rx.Noop[string])

	m.Next("D")

	rx.Pipe1(
		m.Observable,
		rx.DoOnNext(
			func(string) {
				m.Observable.Subscribe(rx.NewBackgroundContext(), rx.Noop[string])
			},
		),
	).Subscribe(rx.NewBackgroundContext(), rx.Noop[string])

	m.Next("E")

	for range 3 {
		m.Observable.Subscribe(rx.NewBackgroundContext(), func(n rx.Notification[string]) {
			panic(ErrTest)
		})
	}

	defer func() {
		if v := recover(); v != ErrTest {
			panic(v)
		}
	}()

	m.Next("F")
}
