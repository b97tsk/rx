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

	NewTestSuite(t).Case(
		d.Observable.Pipe(operators.Take(2)),
		"A", "B", Completed,
	).Case(
		rx.Merge(d.Observable, d.Observable),
		rx.ErrDropped,
	).Case(
		d.Observable,
		"C", Completed,
	).Case(
		d.Observable,
		Completed,
	).TestAll()
}
