package rx_test

import (
	"context"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestObservable_BlockingFirst(t *testing.T) {
	tests := [...]struct {
		obs rx.Observable
		val interface{}
		err error
	}{
		{rx.Empty(), nil, rx.ErrEmpty},
		{rx.Throw(ErrTest), nil, ErrTest},
		{rx.Just("A"), "A", nil},
		{rx.Just("A", "B"), "A", nil},
		{rx.Just("A", "B").Pipe(operators.Go()), "A", nil},
		{rx.Never(), nil, context.DeadlineExceeded},
	}
	ctx, cancel := context.WithTimeout(context.Background(), Step(1))
	for _, x := range tests {
		val, err := x.obs.BlockingFirst(ctx)
		if val != x.val || err != x.err {
			t.Fail()
		}
	}
	cancel()
}

func TestObservable_BlockingFirstOrDefault(t *testing.T) {
	tests := [...]struct {
		obs rx.Observable
		val interface{}
	}{
		{rx.Empty(), "C"},
		{rx.Throw(ErrTest), "C"},
		{rx.Just("A"), "A"},
		{rx.Just("A", "B"), "A"},
		{rx.Just("A", "B").Pipe(operators.Go()), "A"},
		{rx.Never(), "C"},
	}
	ctx, cancel := context.WithTimeout(context.Background(), Step(1))
	for _, x := range tests {
		val := x.obs.BlockingFirstOrDefault(ctx, "C")
		if val != x.val {
			t.Fail()
		}
	}
	cancel()
}

func TestObservable_BlockingLast(t *testing.T) {
	tests := [...]struct {
		obs rx.Observable
		val interface{}
		err error
	}{
		{rx.Empty(), nil, rx.ErrEmpty},
		{rx.Throw(ErrTest), nil, ErrTest},
		{rx.Just("A"), "A", nil},
		{rx.Just("A", "B"), "B", nil},
		{rx.Just("A", "B").Pipe(operators.Go()), "B", nil},
		{rx.Never(), nil, context.DeadlineExceeded},
	}
	ctx, cancel := context.WithTimeout(context.Background(), Step(1))
	for _, x := range tests {
		val, err := x.obs.BlockingLast(ctx)
		if val != x.val || err != x.err {
			t.Fail()
		}
	}
	cancel()
}

func TestObservable_BlockingLastOrDefault(t *testing.T) {
	tests := [...]struct {
		obs rx.Observable
		val interface{}
	}{
		{rx.Empty(), "C"},
		{rx.Throw(ErrTest), "C"},
		{rx.Just("A"), "A"},
		{rx.Just("A", "B"), "B"},
		{rx.Just("A", "B").Pipe(operators.Go()), "B"},
		{rx.Never(), "C"},
	}
	ctx, cancel := context.WithTimeout(context.Background(), Step(1))
	for _, x := range tests {
		val := x.obs.BlockingLastOrDefault(ctx, "C")
		if val != x.val {
			t.Fail()
		}
	}
	cancel()
}

func TestObservable_BlockingSingle(t *testing.T) {
	tests := [...]struct {
		obs rx.Observable
		val interface{}
		err error
	}{
		{rx.Empty(), nil, rx.ErrEmpty},
		{rx.Throw(ErrTest), nil, ErrTest},
		{rx.Just("A"), "A", nil},
		{rx.Just("A", "B"), nil, rx.ErrNotSingle},
		{rx.Just("A", "B").Pipe(operators.Go()), nil, rx.ErrNotSingle},
		{rx.Never(), nil, context.DeadlineExceeded},
	}
	ctx, cancel := context.WithTimeout(context.Background(), Step(1))
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
		obs rx.Observable
		err error
	}{
		{rx.Empty(), nil},
		{rx.Throw(ErrTest), ErrTest},
		{rx.Just("A"), nil},
		{rx.Just("A", "B"), nil},
		{rx.Just("A", "B").Pipe(operators.Go()), nil},
		{rx.Never(), context.DeadlineExceeded},
	}
	ctx, cancel := context.WithTimeout(context.Background(), Step(1))
	for _, x := range tests {
		err := x.obs.BlockingSubscribe(ctx, rx.Noop)
		if err != x.err {
			t.Fail()
		}
	}
	cancel()
}
