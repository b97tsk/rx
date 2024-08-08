package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestBlockingFirst(t *testing.T) {
	t.Parallel()

	testCases := [...]struct {
		ob  rx.Observable[string]
		res rx.Notification[string]
	}{
		{rx.Empty[string](), rx.Error[string](rx.ErrEmpty)},
		{rx.Throw[string](ErrTest), rx.Error[string](ErrTest)},
		{rx.Just("A"), rx.Next("A")},
		{rx.Just("A", "B"), rx.Next("A")},
		{rx.Never[string](), rx.Stop[string](ErrTest)},
	}

	ctx, cancel := rx.NewBackgroundContext().WithTimeoutCause(Step(1), ErrTest)
	defer cancel()

	for _, tc := range testCases {
		if tc.ob.BlockingFirst(ctx) != tc.res {
			t.Fail()
		}
	}
}

func TestBlockingLast(t *testing.T) {
	t.Parallel()

	testCases := [...]struct {
		ob  rx.Observable[string]
		res rx.Notification[string]
	}{
		{rx.Empty[string](), rx.Error[string](rx.ErrEmpty)},
		{rx.Throw[string](ErrTest), rx.Error[string](ErrTest)},
		{rx.Just("A"), rx.Next("A")},
		{rx.Just("A", "B"), rx.Next("B")},
		{rx.Never[string](), rx.Stop[string](ErrTest)},
	}

	ctx, cancel := rx.NewBackgroundContext().WithTimeoutCause(Step(1), ErrTest)
	defer cancel()

	for _, tc := range testCases {
		if tc.ob.BlockingLast(ctx) != tc.res {
			t.Fail()
		}
	}
}

func TestBlockingSingle(t *testing.T) {
	t.Parallel()

	testCases := [...]struct {
		ob  rx.Observable[string]
		res rx.Notification[string]
	}{
		{rx.Empty[string](), rx.Error[string](rx.ErrEmpty)},
		{rx.Throw[string](ErrTest), rx.Error[string](ErrTest)},
		{rx.Just("A"), rx.Next("A")},
		{rx.Just("A", "B"), rx.Error[string](rx.ErrNotSingle)},
		{rx.Never[string](), rx.Stop[string](ErrTest)},
	}

	ctx, cancel := rx.NewBackgroundContext().WithTimeoutCause(Step(1), ErrTest)
	defer cancel()

	for _, tc := range testCases {
		if tc.ob.BlockingSingle(ctx) != tc.res {
			t.Fail()
		}
	}
}

func TestBlockingSubscribe(t *testing.T) {
	t.Parallel()

	testCases := [...]struct {
		ob  rx.Observable[string]
		res rx.Notification[string]
	}{
		{rx.Empty[string](), rx.Complete[string]()},
		{rx.Throw[string](ErrTest), rx.Error[string](ErrTest)},
		{rx.Just("A"), rx.Complete[string]()},
		{rx.Just("A", "B"), rx.Complete[string]()},
		{rx.Never[string](), rx.Stop[string](ErrTest)},
	}

	ctx, cancel := rx.NewBackgroundContext().WithTimeoutCause(Step(1), ErrTest)
	defer cancel()

	for _, tc := range testCases {
		if tc.ob.BlockingSubscribe(ctx, rx.Noop[string]) != tc.res {
			t.Fail()
		}
	}
}
