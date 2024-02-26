package rx_test

import (
	"context"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestBlockingFirst(t *testing.T) {
	t.Parallel()

	testCases := [...]struct {
		obs rx.Observable[string]
		val string
		err error
	}{
		{rx.Empty[string](), "", rx.ErrEmpty},
		{rx.Throw[string](ErrTest), "", ErrTest},
		{rx.Just("A"), "A", nil},
		{rx.Just("A", "B"), "A", nil},
		{rx.Never[string](), "", context.DeadlineExceeded},
	}

	ctx, cancel := rx.NewBackgroundContext().WithTimeout(Step(1))
	defer cancel()

	for _, tc := range testCases {
		v, err := tc.obs.BlockingFirst(ctx)
		if v != tc.val || err != tc.err {
			t.Fail()
		}
	}
}

func TestBlockingFirstOrElse(t *testing.T) {
	t.Parallel()

	testCases := [...]struct {
		obs rx.Observable[string]
		val string
	}{
		{rx.Empty[string](), "C"},
		{rx.Throw[string](ErrTest), "C"},
		{rx.Just("A"), "A"},
		{rx.Just("A", "B"), "A"},
		{rx.Never[string](), "C"},
	}

	ctx, cancel := rx.NewBackgroundContext().WithTimeout(Step(1))
	defer cancel()

	for _, tc := range testCases {
		v := tc.obs.BlockingFirstOrElse(ctx, "C")
		if v != tc.val {
			t.Fail()
		}
	}
}

func TestBlockingLast(t *testing.T) {
	t.Parallel()

	testCases := [...]struct {
		obs rx.Observable[string]
		val string
		err error
	}{
		{rx.Empty[string](), "", rx.ErrEmpty},
		{rx.Throw[string](ErrTest), "", ErrTest},
		{rx.Just("A"), "A", nil},
		{rx.Just("A", "B"), "B", nil},
		{rx.Never[string](), "", context.DeadlineExceeded},
	}

	ctx, cancel := rx.NewBackgroundContext().WithTimeout(Step(1))
	defer cancel()

	for _, tc := range testCases {
		v, err := tc.obs.BlockingLast(ctx)
		if v != tc.val || err != tc.err {
			t.Fail()
		}
	}
}

func TestBlockingLastOrElse(t *testing.T) {
	t.Parallel()

	testCases := [...]struct {
		obs rx.Observable[string]
		val string
	}{
		{rx.Empty[string](), "C"},
		{rx.Throw[string](ErrTest), "C"},
		{rx.Just("A"), "A"},
		{rx.Just("A", "B"), "B"},
		{rx.Never[string](), "C"},
	}

	ctx, cancel := rx.NewBackgroundContext().WithTimeout(Step(1))
	defer cancel()

	for _, tc := range testCases {
		v := tc.obs.BlockingLastOrElse(ctx, "C")
		if v != tc.val {
			t.Fail()
		}
	}
}

func TestBlockingSingle(t *testing.T) {
	t.Parallel()

	testCases := [...]struct {
		obs rx.Observable[string]
		val string
		err error
	}{
		{rx.Empty[string](), "", rx.ErrEmpty},
		{rx.Throw[string](ErrTest), "", ErrTest},
		{rx.Just("A"), "A", nil},
		{rx.Just("A", "B"), "", rx.ErrNotSingle},
		{rx.Never[string](), "", context.DeadlineExceeded},
	}

	ctx, cancel := rx.NewBackgroundContext().WithTimeout(Step(1))
	defer cancel()

	for _, tc := range testCases {
		v, err := tc.obs.BlockingSingle(ctx)
		if v != tc.val || err != tc.err {
			t.Fail()
		}
	}
}

func TestBlockingSubscribe(t *testing.T) {
	t.Parallel()

	testCases := [...]struct {
		obs rx.Observable[string]
		err error
	}{
		{rx.Empty[string](), nil},
		{rx.Throw[string](ErrTest), ErrTest},
		{rx.Just("A"), nil},
		{rx.Just("A", "B"), nil},
		{rx.Never[string](), context.DeadlineExceeded},
	}

	ctx, cancel := rx.NewBackgroundContext().WithTimeout(Step(1))
	defer cancel()

	for _, tc := range testCases {
		err := tc.obs.BlockingSubscribe(ctx, rx.Noop[string])
		if err != tc.err {
			t.Fail()
		}
	}
}
