package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestFromChannel(t *testing.T) {
	t.Run("#1", func(t *testing.T) {
		c := make(chan interface{})
		go func() {
			c <- "A"
			c <- "B"
			c <- "C"
			close(c)
		}()
		Subscribe(t, rx.FromChannel(c), "A", "B", "C", rx.Completed)
	})
	t.Run("#2", func(t *testing.T) {
		c := make(chan interface{})
		go func() {
			c <- "A"
			c <- "B"
			c <- "C"
			close(c)
		}()
		SubscribeN(
			t,
			[]rx.Observable{
				rx.FromChannel(c).Pipe(operators.Take(1)),
				rx.FromChannel(c).Pipe(operators.Take(2)),
			},
			[][]interface{}{
				{"A", rx.Completed},
				{"B", "C", rx.Completed},
			},
		)
	})
}
