package rx_test

import (
	"context"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestObservable_Publish(t *testing.T) {
	obs := rx.Ticker(Step(1)).Pipe(
		operators.Map(
			func(val interface{}, idx int) interface{} {
				return idx
			},
		),
	).Publish()
	ctx, _ := rx.Zip(
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
		ToString(),
	).Subscribe(context.Background(), func(x rx.Notification) {
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
	_, disconnect := obs.Connect(context.Background())
	<-ctx.Done()
	disconnect()
}

func TestObservable_PublishBehavior(t *testing.T) {
	obs := rx.Ticker(Step(1)).Pipe(
		operators.Map(
			func(val interface{}, idx int) interface{} {
				return idx
			},
		),
	).PublishBehavior(-1)
	ctx, _ := rx.Zip(
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
		ToString(),
	).Subscribe(context.Background(), func(x rx.Notification) {
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
	_, disconnect := obs.Connect(context.Background())
	<-ctx.Done()
	disconnect()
}

func TestObservable_PublishReplay(t *testing.T) {
	obs := rx.Ticker(Step(2)).Pipe(
		operators.Map(
			func(val interface{}, idx int) interface{} {
				return idx
			},
		),
	).PublishReplay(2, 0)
	ctx, _ := rx.Zip(
		obs.Pipe(operators.Take(4)),
		obs.Pipe(operators.Skip(4), operators.Take(4), DelaySubscription(7)),
	).Pipe(
		operators.Map(
			func(val interface{}, idx int) interface{} {
				vals := val.([]interface{})
				return vals[0].(int) * vals[1].(int)
			},
		),
		operators.ToSlice(),
		ToString(),
	).Subscribe(context.Background(), func(x rx.Notification) {
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
	_, disconnect := obs.Connect(context.Background())
	<-ctx.Done()
	disconnect()
}
