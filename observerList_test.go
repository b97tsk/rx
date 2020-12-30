package rx_test

import (
	"context"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestObserverList(t *testing.T) {
	d := rx.Multicast()

	rx.Just("A", "B", "C").Pipe(
		AddLatencyToValues(1, 1),
	).Subscribe(context.Background(), d.Observer.ElementsOnly)

	Test(t, d.Observable.Pipe(operators.Take(2)), "A", "B", Completed)

	time.Sleep(time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	d.Observable.Pipe(
		operators.DoOnNext(
			func(interface{}) {
				time.Sleep(time.Second)
			},
		),
	).Subscribe(ctx, rx.Noop)

	d.Observer.Next("D")

	d.Observable.Pipe(
		operators.DoOnNext(
			func(interface{}) {
				d.Observable.Subscribe(context.Background(), rx.Noop)
			},
		),
	).Subscribe(context.Background(), rx.Noop)

	d.Observer.Next("E")
	d.Observer.Complete()
}
