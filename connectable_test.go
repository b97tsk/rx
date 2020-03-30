package rx_test

import (
	"context"
	"testing"

	. "github.com/b97tsk/rx"
)

func TestObservable_Publish(t *testing.T) {
	obs := Interval(step(1)).Publish()
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
	).Subscribe(context.Background(), func(x Notification) {
		switch {
		case x.HasValue:
			if x.Value != "[0 5 12 21]" {
				t.Fail()
			}
		case x.HasError:
			t.Error(x.Error)
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

func TestObservable_PublishBehavior(t *testing.T) {
	obs := Interval(step(1)).PublishBehavior(-1)
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
	).Subscribe(context.Background(), func(x Notification) {
		switch {
		case x.HasValue:
			if x.Value != "[-3 0 5 12]" {
				t.Fail()
			}
		case x.HasError:
			t.Error(x.Error)
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

func TestObservable_PublishReplay(t *testing.T) {
	obs := Interval(step(2)).PublishReplay(2, 0)
	ctx, _ := Zip(
		obs.Pipe(operators.Take(4)),
		obs.Pipe(operators.Skip(4), operators.Take(4), delaySubscription(7)),
	).Pipe(
		operators.Map(
			func(val interface{}, idx int) interface{} {
				vals := val.([]interface{})
				return vals[0].(int) * vals[1].(int)
			},
		),
		operators.ToSlice(),
		toString,
	).Subscribe(context.Background(), func(x Notification) {
		switch {
		case x.HasValue:
			if x.Value != "[0 6 14 24]" {
				t.Fail()
			}
		case x.HasError:
			t.Error(x.Error)
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
