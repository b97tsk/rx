package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestFromChan1(t *testing.T) {
	c := make(chan interface{})

	go func() {
		c <- "A"
		c <- "B"
		c <- "C"
		close(c)
	}()

	NewTestSuite(t).Case(
		rx.FromChan(c),
		"A", "B", "C", Completed,
	).TestAll()
}

func TestFromChan2(t *testing.T) {
	c := make(chan interface{})

	go func() {
		c <- "A"
		c <- "B"
		c <- "C"
		close(c)
	}()

	obs := rx.FromChan(c)

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
}
