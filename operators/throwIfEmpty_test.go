package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestThrowIfEmpty(t *testing.T) {
	throwIfEmpty := operators.ThrowIfEmpty(rx.ErrEmpty)
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Empty().Pipe(throwIfEmpty),
			rx.Throw(ErrTest).Pipe(throwIfEmpty),
			rx.Just(1).Pipe(throwIfEmpty),
			rx.Just(1, 2).Pipe(throwIfEmpty),
			rx.Concat(rx.Just(1), rx.Throw(ErrTest)).Pipe(throwIfEmpty),
			rx.Concat(rx.Just(1, 2), rx.Throw(ErrTest)).Pipe(throwIfEmpty),
		},
		[][]interface{}{
			{rx.ErrEmpty},
			{ErrTest},
			{1, rx.Completed},
			{1, 2, rx.Completed},
			{1, ErrTest},
			{1, 2, ErrTest},
		},
	)
}
