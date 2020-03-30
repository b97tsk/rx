package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestRepeatWhen(t *testing.T) {
	var (
		repeatOnce  = operators.Take(0)
		repeatTwice = operators.Take(1)
	)

	SubscribeN(
		t,
		[]rx.Observable{
			rx.Range(1, 4).Pipe(operators.RepeatWhen(repeatOnce)),
			rx.Range(1, 4).Pipe(operators.RepeatWhen(repeatTwice)),
			rx.Concat(rx.Range(1, 4), rx.Throw(ErrTest)).Pipe(operators.RepeatWhen(repeatOnce)),
			rx.Concat(rx.Range(1, 4), rx.Throw(ErrTest)).Pipe(operators.RepeatWhen(repeatTwice)),
		},
		[][]interface{}{
			{1, 2, 3, rx.Complete},
			{1, 2, 3, 1, 2, 3, rx.Complete},
			{1, 2, 3, ErrTest},
			{1, 2, 3, ErrTest},
		},
	)
}
