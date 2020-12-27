package rx_test

import (
	"context"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestFromChan(t *testing.T) {
	makeChan := func() <-chan interface{} {
		c := make(chan interface{})

		go func() {
			c <- "A"
			c <- "B"
			c <- "C"
			close(c)
		}()

		return c
	}

	NewTestSuite(t).Case(
		rx.FromChan(makeChan()),
		"A", "B", "C", Completed,
	).TestAll()

	obs := rx.FromChan(makeChan())

	NewTestSuite(t).Case(
		obs.Pipe(operators.Take(1)),
		"A", Completed,
	).Case(
		obs.Pipe(operators.Take(2)),
		"B", "C", Completed,
	).Case(
		obs,
		Completed,
	).TestAll()

	ctx, cancel := context.WithTimeout(context.Background(), Step(1))
	defer cancel()

	rx.FromChan(nil).Subscribe(ctx, rx.Noop)
}
