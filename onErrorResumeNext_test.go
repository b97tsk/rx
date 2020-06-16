package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestOnErrorResumeNext(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.OnErrorResumeNext(rx.Range(1, 4), rx.Range(4, 7), rx.Range(7, 10)),
			rx.OnErrorResumeNext(rx.Throw(ErrTest), rx.Range(4, 7), rx.Range(7, 10)),
			rx.OnErrorResumeNext(rx.Range(1, 4), rx.Throw(ErrTest), rx.Range(7, 10)),
			rx.OnErrorResumeNext(rx.Range(1, 4), rx.Range(4, 7), rx.Throw(ErrTest)),
		},
		[][]interface{}{
			{1, 2, 3, 4, 5, 6, 7, 8, 9, rx.Completed},
			{4, 5, 6, 7, 8, 9, rx.Completed},
			{1, 2, 3, 7, 8, 9, rx.Completed},
			{1, 2, 3, 4, 5, 6, rx.Completed},
		},
	)
}
