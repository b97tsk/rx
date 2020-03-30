package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestEvery(t *testing.T) {
	everyLessThan5 := operators.Every(
		func(val interface{}, idx int) bool {
			return val.(int) < 5
		},
	)
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Range(1, 9).Pipe(everyLessThan5),
			rx.Range(1, 5).Pipe(everyLessThan5),
			rx.Empty().Pipe(everyLessThan5),
			rx.Concat(rx.Range(1, 9), rx.Throw(ErrTest)).Pipe(everyLessThan5),
			rx.Concat(rx.Range(1, 5), rx.Throw(ErrTest)).Pipe(everyLessThan5),
		},
		[][]interface{}{
			{false, rx.Complete},
			{true, rx.Complete},
			{true, rx.Complete},
			{false, rx.Complete},
			{ErrTest},
		},
	)
}
