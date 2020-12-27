package rx_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestRace(t *testing.T) {
	const N = 8

	var observables [N]rx.Observable

	for i := 0; i < N; i++ {
		d := Step(i + 1)
		observables[i] = rx.Timer(d).Pipe(operators.MapTo(d))
	}

	rand.Seed(time.Now().UnixNano())

	rand.Shuffle(N, func(i, j int) {
		observables[i], observables[j] = observables[j], observables[i]
	})

	NewTestSuite(t).Case(
		rx.Race(observables[:]...),
		Step(1),
		Completed,
	).Case(
		rx.Race(),
		Completed,
	).TestAll()
}
