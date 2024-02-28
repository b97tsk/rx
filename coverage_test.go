package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestAdditionalCoverage(t *testing.T) {
	t.Parallel()

	t.Run("ComposeFuncs", func(t *testing.T) {
		obs := rx.Empty[int]()
		op := rx.SkipAll[int]()
		_ = rx.Compose2(op, op).Apply(obs)
		_ = rx.Compose3(op, op, op).Apply(obs)
		_ = rx.Compose4(op, op, op, op).Apply(obs)
		_ = rx.Compose5(op, op, op, op, op).Apply(obs)
		_ = rx.Compose6(op, op, op, op, op, op).Apply(obs)
		_ = rx.Compose7(op, op, op, op, op, op, op).Apply(obs)
		_ = rx.Compose8(op, op, op, op, op, op, op, op).Apply(obs)
		_ = rx.Compose9(op, op, op, op, op, op, op, op, op).Apply(obs)
	})

	t.Run("NilObservable", func(t *testing.T) {
		_ = rx.NewObservable[any](nil).BlockingSubscribe(rx.NewBackgroundContext(), func(n rx.Notification[any]) {
			if n.Kind != rx.KindError || n.Error != rx.ErrNil {
				t.Fail()
			}
		})
	})

	t.Run("NewObserver", func(t *testing.T) {
		observer := rx.NewObserver(rx.Noop[int])
		observer.Emit(rx.Next(0))
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
