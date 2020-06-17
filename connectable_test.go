package rx_test

import (
	"context"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestObservable_Publish(t *testing.T) {
	t.Run("#1", observablePublishTest1)
	t.Run("#2", observablePublishTest2)
	t.Run("#3", observablePublishTest3)
}

func observablePublishTest1(t *testing.T) {
	obs := rx.Ticker(Step(1)).Pipe(
		operators.Map(
			func(val interface{}, idx int) interface{} {
				return idx
			},
		),
	).Publish(rx.NewSubject().Double)
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
	).Subscribe(context.Background(), func(u rx.Notification) {
		switch {
		case u.HasValue:
			if u.Value != "[0 5 12 21]" {
				t.Fail()
			}
		case u.HasError:
			t.Error(u.Error)
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

func observablePublishTest2(t *testing.T) {
	obs := rx.Ticker(Step(1)).Pipe(
		operators.Map(
			func(val interface{}, idx int) interface{} {
				return idx
			},
		),
	).Publish(rx.NewBehaviorSubject(-1).Double)
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
	).Subscribe(context.Background(), func(u rx.Notification) {
		switch {
		case u.HasValue:
			if u.Value != "[-3 0 5 12]" {
				t.Fail()
			}
		case u.HasError:
			t.Error(u.Error)
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

func observablePublishTest3(t *testing.T) {
	obs := rx.Ticker(Step(2)).Pipe(
		operators.Map(
			func(val interface{}, idx int) interface{} {
				return idx
			},
		),
	).Publish(rx.NewReplaySubject(2).Double)
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
	).Subscribe(context.Background(), func(u rx.Notification) {
		switch {
		case u.HasValue:
			if u.Value != "[0 6 14 24]" {
				t.Fail()
			}
		case u.HasError:
			t.Error(u.Error)
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
