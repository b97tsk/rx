package rx_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	. "github.com/b97tsk/rx"
)

func TestRace(t *testing.T) {
	const N = 8

	var observables [N]Observable

	for i := 0; i < N; i++ {
		d := step(i + 1)
		observables[i] = Just(d).Pipe(operators.Delay(d))
	}

	rand.Seed(time.Now().UnixNano())

	rand.Shuffle(N, func(i, j int) {
		observables[i], observables[j] = observables[j], observables[i]
	})

	winner, err := Race(observables[:]...).BlockingSingle(context.Background())
	if err != nil {
		t.Error(err)
	}

	if winner != step(1) {
		t.Log("Winner:", winner)
		t.Fail()
	}
}
