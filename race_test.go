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

	var s [8]rx.Observable[string]

	detachContext := rx.NewOperator(
		func(source rx.Observable[string]) rx.Observable[string] {
			return func(_ rx.Context, o rx.Observer[string]) {
				source.Subscribe(rx.NewBackgroundContext(), o)
			}
		},
	)

	for i := range s {
		s[i] = rx.Pipe2(
			rx.Timer(Step(i+1)),
			rx.MapTo[time.Time](string(rune(i+'A'))),
			detachContext,
		)
	}

	rand.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })

	NewTestSuite[string](t).Case(
		rx.Race[string](),
		ErrComplete,
	).Case(
		rx.Race(s[:]...),
		"A", ErrComplete,
	).Case(
		rx.Pipe1(s[0], rx.RaceWith(s[1:]...)),
		"A", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Empty[string](),
			rx.RaceWith(s[:]...),
		),
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Race(s[:len(s)-1]...),
			rx.RaceWith(s[len(s)-1], rx.Just("A"), rx.Just("B"), rx.Just("C")),
		),
		"A", ErrComplete,
	)

	time.Sleep(Step(len(s)))
}
