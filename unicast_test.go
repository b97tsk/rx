package rx_test

import (
	"context"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestUnicast(t *testing.T) {
	d := rx.Unicast()

	rx.Just("A", "B", "C").Pipe(
		AddLatencyToNotifications(1, 1),
	).Subscribe(context.Background(), d.Observer)

	NewTestSuite(t).Case(
		d.Observable.Pipe(operators.Take(2)),
		"A", "B", Completed,
	).Case(
		d.Observable,
		rx.ErrDropped,
	).Case(
		d.Observable.Pipe(DelaySubscription(3)),
		Completed,
	).Case(
		d.Observable,
		Completed,
	).TestAll()
}
