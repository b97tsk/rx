package rx_test

import (
	"context"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestUnicastLatest(t *testing.T) {
	d := rx.UnicastLatest()

	rx.Just("A", "B", "C").Pipe(
		AddLatencyToNotifications(1, 1),
	).Subscribe(context.Background(), d.Observer)

	Subscribe(t, d.Observable.Pipe(operators.Take(2)), "A", "B", Completed)
	Subscribe(t, rx.Merge(d.Observable, d.Observable), rx.ErrDropped)
	Subscribe(t, d.Observable, "C", Completed)
	Subscribe(t, d.Observable, Completed)
}
