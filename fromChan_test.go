package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestFromChan(t *testing.T) {
	t.Run("#1", func(t *testing.T) {
		c := make(chan interface{})
		go func() {
			c <- "A"
			c <- "B"
			c <- "C"
			close(c)
		}()
		Subscribe(t, rx.FromChan(c), "A", "B", "C", rx.Completed)
	})
	t.Run("#2", func(t *testing.T) {
		c := make(chan interface{})
		go func() {
			c <- "A"
			c <- "B"
			c <- "C"
			close(c)
		}()
		obs := rx.FromChan(c)
		SubscribeN(
			t,
			[]rx.Observable{
				obs.Pipe(operators.Take(1)),
				obs.Pipe(operators.Take(2)),
				obs,
			},
			[][]interface{}{
				{"A", rx.Completed},
				{"B", "C", rx.Completed},
				{rx.Completed},
			},
		)
	})
}
