package rx_test

import (
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestAdditionalCoverage(t *testing.T) {
	t.Parallel()

	t.Run("ComposeFuncs", func(t *testing.T) {
		ob := rx.Empty[int]()
		op := rx.Discard[int]()
		_ = rx.Compose2(op, op).Apply(ob)
		_ = rx.Compose3(op, op, op).Apply(ob)
		_ = rx.Compose4(op, op, op, op).Apply(ob)
		_ = rx.Compose5(op, op, op, op, op).Apply(ob)
		_ = rx.Compose6(op, op, op, op, op, op).Apply(ob)
		_ = rx.Compose7(op, op, op, op, op, op, op).Apply(ob)
		_ = rx.Compose8(op, op, op, op, op, op, op, op).Apply(ob)
		_ = rx.Compose9(op, op, op, op, op, op, op, op, op).Apply(ob)
	})

	t.Run("Context.WithCancelCause", func(t *testing.T) {
		ctx, cancel := rx.NewBackgroundContext().WithCancelCause()
		cancel(ErrTest)
		if ctx.Cause() != ErrTest {
			t.Fail()
		}
	})

	t.Run("Context.WithDeadlineCause", func(t *testing.T) {
		ctx, cancel := rx.NewBackgroundContext().WithDeadlineCause(time.Now().Add(Step(1)), ErrTest)
		defer cancel()
		time.Sleep(Step(2))
		if ctx.Cause() != ErrTest {
			t.Fail()
		}
	})

	t.Run("Context.WithTimeoutCause", func(t *testing.T) {
		ctx, cancel := rx.NewBackgroundContext().WithTimeoutCause(Step(1), ErrTest)
		defer cancel()
		time.Sleep(Step(2))
		if ctx.Cause() != ErrTest {
			t.Fail()
		}
	})

	t.Run("Go", func(t *testing.T) {
		n := rx.Pipe1(rx.Just(42), rx.Go[int]()).BlockingSingle(rx.NewBackgroundContext())
		if n.Error != nil {
			t.Log(n.Error)
		}
		if n.Value != 42 {
			t.Fail()
		}
	})

	t.Run("Notification.And", func(t *testing.T) {
		if rx.Next(42).And(rx.Error[int](ErrTest)).Error != ErrTest {
			t.Fail()
		}
		if rx.Error[int](ErrTest).And(rx.Next(42)).Error != ErrTest {
			t.Fail()
		}
	})

	t.Run("Notification.Or", func(t *testing.T) {
		if rx.Error[int](ErrTest).Or(rx.Next(42)).Value != 42 {
			t.Fail()
		}
		if rx.Next(42).Or(rx.Error[int](ErrTest)).Value != 42 {
			t.Fail()
		}
	})

	t.Run("NilObservable", func(t *testing.T) {
		NewTestSuite[int](t).Case(rx.NewObservable[int](nil), rx.ErrOops, "nil Observable")
	})

	t.Run("NewObserver", func(t *testing.T) {
		o := rx.NewObserver(rx.Noop[int])
		o.Unsubscribe()
	})

	t.Run("TryFuncs", func(t *testing.T) {
		defer func() {
			if v := recover(); v != ErrTest {
				panic(v)
			}
		}()

		defer rx.Try0(func() { panic(ErrTest) }, func() {})
		defer rx.Try1(func(any) { panic(ErrTest) }, nil, func() {})
		defer rx.Try2(func(any, any) { panic(ErrTest) }, nil, nil, func() {})
		defer rx.Try3(func(any, any, any) { panic(ErrTest) }, nil, nil, nil, func() {})
		defer rx.Try4(func(any, any, any, any) { panic(ErrTest) }, nil, nil, nil, nil, func() {})
		defer rx.Try5(func(any, any, any, any, any) { panic(ErrTest) }, nil, nil, nil, nil, nil, func() {})
		defer rx.Try6(func(any, any, any, any, any, any) { panic(ErrTest) }, nil, nil, nil, nil, nil, nil, func() {})
		defer rx.Try7(func(any, any, any, any, any, any, any) { panic(ErrTest) }, nil, nil, nil, nil, nil, nil, nil, func() {})
		defer rx.Try8(func(any, any, any, any, any, any, any, any) { panic(ErrTest) }, nil, nil, nil, nil, nil, nil, nil, nil, func() {})
		defer rx.Try9(func(any, any, any, any, any, any, any, any, any) { panic(ErrTest) }, nil, nil, nil, nil, nil, nil, nil, nil, nil, func() {})

		defer rx.Try01(func() any { panic(ErrTest) }, func() {})
		defer rx.Try11(func(any) any { panic(ErrTest) }, nil, func() {})
		defer rx.Try21(func(any, any) any { panic(ErrTest) }, nil, nil, func() {})
		defer rx.Try31(func(any, any, any) any { panic(ErrTest) }, nil, nil, nil, func() {})
		defer rx.Try41(func(any, any, any, any) any { panic(ErrTest) }, nil, nil, nil, nil, func() {})
		defer rx.Try51(func(any, any, any, any, any) any { panic(ErrTest) }, nil, nil, nil, nil, nil, func() {})
		defer rx.Try61(func(any, any, any, any, any, any) any { panic(ErrTest) }, nil, nil, nil, nil, nil, nil, func() {})
		defer rx.Try71(func(any, any, any, any, any, any, any) any { panic(ErrTest) }, nil, nil, nil, nil, nil, nil, nil, func() {})
		defer rx.Try81(func(any, any, any, any, any, any, any, any) any { panic(ErrTest) }, nil, nil, nil, nil, nil, nil, nil, nil, func() {})
		defer rx.Try91(func(any, any, any, any, any, any, any, any, any) any { panic(ErrTest) }, nil, nil, nil, nil, nil, nil, nil, nil, nil, func() {})
	})
}
