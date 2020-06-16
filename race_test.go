package rx_test

import (
	"context"
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

	winner, err := rx.Race(observables[:]...).BlockingSingle(context.Background())
	if err != nil {
		t.Error(err)
	}

	if winner != Step(1) {
		t.Log("Winner:", winner)
		t.Fail()
	}
}
