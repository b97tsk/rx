package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestChannelize(t *testing.T) {
	t.Parallel()

	join := func(upstream <-chan rx.Notification[string], downstream chan<- rx.Notification[string]) {
		for n := range upstream {
			switch n.Kind {
			case rx.KindNext:
				downstream <- n
			case rx.KindError, rx.KindComplete:
				downstream <- n
				return
			}
		}
	}

	NewTestSuite[string](t).Case(
		rx.Pipe1(
			rx.Just("A", "B", "C"),
			rx.Channelize(join),
		),
		"A", "B", "C", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Just("A", "B", "C"),
			rx.Channelize(
				func(upstream <-chan rx.Notification[string], downstream chan<- rx.Notification[string]) {
					panic(ErrTest)
				},
			),
		),
		rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			rx.Channelize(join),
			rx.OnNext(func(string) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}
