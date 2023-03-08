package rx_test

import (
	"context"
	"testing"

	"github.com/b97tsk/rx"
)

func TestAdditionalCoverage(t *testing.T) {
	t.Parallel()

	observer := rx.NewObserver(rx.Noop[int])
	observer.Emit(rx.Next(0))

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

	_ = rx.NewObservable[any](nil).BlockingSubscribe(
		context.Background(),
		func(n rx.Notification[any]) {
			if !n.HasError || n.Error != rx.ErrNil {
				t.Fail()
			}
		},
	)
}
