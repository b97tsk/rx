package rx_test

import (
	"context"
	"testing"

	. "github.com/b97tsk/rx"
)

func TestObservable_BlockingFirst(t *testing.T) {
	tests := [...]struct {
		obs Observable
		val interface{}
		err error
	}{
		{Empty(), nil, ErrEmpty},
		{Throw(xErrTest), nil, xErrTest},
		{Just("A"), "A", nil},
		{Just("A", "B"), "A", nil},
		{Just("A", "B").Pipe(operators.Go()), "A", nil},
		{Never(), nil, context.DeadlineExceeded},
	}
	ctx, cancel := context.WithTimeout(context.Background(), step(1))
	for _, x := range tests {
		val, err := x.obs.BlockingFirst(ctx)
		if val != x.val || err != x.err {
			t.Fail()
		}
	}
	cancel()
}

func TestObservable_BlockingLast(t *testing.T) {
	tests := [...]struct {
		obs Observable
		val interface{}
		err error
	}{
		{Empty(), nil, ErrEmpty},
		{Throw(xErrTest), nil, xErrTest},
		{Just("A"), "A", nil},
		{Just("A", "B"), "B", nil},
		{Just("A", "B").Pipe(operators.Go()), "B", nil},
		{Never(), nil, context.DeadlineExceeded},
	}
	ctx, cancel := context.WithTimeout(context.Background(), step(1))
	for _, x := range tests {
		val, err := x.obs.BlockingLast(ctx)
		if val != x.val || err != x.err {
			t.Fail()
		}
	}
	cancel()
}

func TestObservable_BlockingSingle(t *testing.T) {
	tests := [...]struct {
		obs Observable
		val interface{}
		err error
	}{
		{Empty(), nil, ErrEmpty},
		{Throw(xErrTest), nil, xErrTest},
		{Just("A"), "A", nil},
		{Just("A", "B"), nil, ErrNotSingle},
		{Just("A", "B").Pipe(operators.Go()), nil, ErrNotSingle},
		{Never(), nil, context.DeadlineExceeded},
	}
	ctx, cancel := context.WithTimeout(context.Background(), step(1))
	for _, x := range tests {
		val, err := x.obs.BlockingSingle(ctx)
		if val != x.val || err != x.err {
			t.Fail()
		}
	}
	cancel()
}

func TestObservable_BlockingSubscribe(t *testing.T) {
	tests := [...]struct {
		obs Observable
		err error
	}{
		{Empty(), nil},
		{Throw(xErrTest), xErrTest},
		{Just("A"), nil},
		{Just("A", "B"), nil},
		{Just("A", "B").Pipe(operators.Go()), nil},
		{Never(), context.DeadlineExceeded},
	}
	ctx, cancel := context.WithTimeout(context.Background(), step(1))
	for _, x := range tests {
		err := x.obs.BlockingSubscribe(ctx, NopObserver)
		if err != x.err {
			t.Fail()
		}
	}
	cancel()
}
