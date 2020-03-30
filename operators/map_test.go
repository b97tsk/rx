package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestMap(t *testing.T) {
	op := operators.Map(
		func(val interface{}, idx int) interface{} {
			return val.(int) * 2
		},
	)
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Empty().Pipe(op),
			rx.Range(1, 5).Pipe(op),
			rx.Concat(rx.Range(1, 5), rx.Throw(ErrTest)).Pipe(op),
		},
		[][]interface{}{
			{rx.Complete},
			{2, 4, 6, 8, rx.Complete},
			{2, 4, 6, 8, ErrTest},
		},
	)
}

func TestMapTo(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Empty().Pipe(operators.MapTo(42)),
			rx.Just("A", "B", "C").Pipe(operators.MapTo(42)),
			rx.Concat(rx.Just("A", "B", "C"), rx.Throw(ErrTest)).Pipe(operators.MapTo(42)),
		},
		[][]interface{}{
			{rx.Complete},
			{42, 42, 42, rx.Complete},
			{42, 42, 42, ErrTest},
		},
	)
}
