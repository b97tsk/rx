package rx_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestRace(t *testing.T) {
	t.Parallel()

	var s [8]rx.Observable[time.Duration]

	detachContext := rx.NewOperator(
		func(source rx.Observable[time.Duration]) rx.Observable[time.Duration] {
			return func(_ rx.Context, sink rx.Observer[time.Duration]) {
				source.Subscribe(rx.NewBackgroundContext(), sink)
			}
		},
	)

	for i := range s {
		d := Step(i + 1)
		s[i] = rx.Pipe2(
			rx.Timer(d),
			rx.MapTo[time.Time](d),
			detachContext,
		)
	}

	rand.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })

	NewTestSuite[time.Duration](t).Case(
		rx.Race[time.Duration](),
		ErrComplete,
	).Case(
		rx.Race(s[:]...),
		Step(1),
		ErrComplete,
	).Case(
		rx.Pipe1(s[0], rx.RaceWith(s[1:]...)),
		Step(1),
		ErrComplete,
	)

	time.Sleep(Step(len(s)))
}
