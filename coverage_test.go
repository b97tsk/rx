package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
)

func TestAdditionalCoverage(t *testing.T) {
	t.Parallel()

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

	_ = rx.NewObservable[any](nil).BlockingSubscribe(rx.NewBackgroundContext(), func(n rx.Notification[any]) {
		if n.Kind != rx.KindError || n.Error != rx.ErrNil {
			t.Fail()
		}
	})

	observer := rx.NewObserver(rx.Noop[int])
	observer.Emit(rx.Next(0))
}
