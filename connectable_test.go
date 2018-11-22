package rx_test

import (
	"context"
	"testing"

	. "github.com/b97tsk/rx"
)

func TestConnectableObservable(t *testing.T) {
	obs := Interval(step(3)).Publish()
	ctx, _ := Zip(
		obs.Pipe(operators.Take(4)),
		obs.Pipe(operators.Skip(4), operators.Take(4)),
	).Pipe(
		operators.Map(
			func(val interface{}, idx int) interface{} {
				vals := val.([]interface{})
				return vals[0].(int) * vals[1].(int)
			},
		),
		operators.ToSlice(),
		toString,
	).Subscribe(context.Background(), func(tt Notification) {
		switch {
		case tt.HasValue:
			if tt.Value != "[0 5 12 21]" {
				t.Fail()
			}
		case tt.HasError:
			t.Error(tt.Value)
		}
	})
	select {
	case <-ctx.Done():
		t.Fail()
		return
	default:
	}
	_, disconnect := obs.Connect()
	defer disconnect()
	<-ctx.Done()
}

func TestOperators_Share(t *testing.T) {
	t.Run("#1", func(t *testing.T) {
		obs := Interval(step(3)).Pipe(
			operators.Take(4),
			operators.Share(),
		)
		subscribe(
			t,
			[]Observable{
				Merge(
					obs,
					obs.Pipe(delaySubscription(4)),
					obs.Pipe(delaySubscription(8)),
					obs.Pipe(delaySubscription(15)),
				),
			},
			0, 1, 1, 2, 2, 2, 3, 3, 3, 0, 1, 2, 3, xComplete,
		)
	})
	t.Run("#2", func(t *testing.T) {
		obs := Interval(step(3)).Pipe(
			operators.Share(),
			operators.Take(4),
		)
		subscribe(
			t,
			[]Observable{
				Merge(
					obs,
					obs.Pipe(delaySubscription(4)),
					obs.Pipe(delaySubscription(8)),
					obs.Pipe(delaySubscription(20)),
				),
			},
			0, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 5, 0, 1, 2, 3, xComplete,
		)
	})
}
