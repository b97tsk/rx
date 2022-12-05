package rx_test

import (
	"context"
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
			return func(ctx context.Context, sink rx.Observer[time.Duration]) {
				source.Subscribe(context.Background(), sink)
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

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })

	NewTestSuite[time.Duration](t).Case(
		rx.Race[time.Duration](),
		ErrCompleted,
	).Case(
		rx.Race(s[:]...),
		Step(1),
		ErrCompleted,
	).Case(
		rx.Pipe(s[0], rx.RaceWith(s[1:]...)),
		Step(1),
		ErrCompleted,
	)

	time.Sleep(Step(len(s)))
}
