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
	}
	for _, x := range tests {
		val, err := x.obs.BlockingFirst(context.Background())
		if val != x.val || err != x.err {
			t.Fail()
		}
	}
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
	}
	for _, x := range tests {
		val, err := x.obs.BlockingLast(context.Background())
		if val != x.val || err != x.err {
			t.Fail()
		}
	}
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
		{Just("A", "B"), "A", ErrNotSingle},
		{Just("A", "B").Pipe(operators.Go()), "A", ErrNotSingle},
	}
	for _, x := range tests {
		val, err := x.obs.BlockingSingle(context.Background())
		if val != x.val || err != x.err {
			t.Fail()
		}
	}
}
