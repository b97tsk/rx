package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOnErrorResumeNext(t *testing.T) {
	subscribeN(
		t,
		[]Observable{
			OnErrorResumeNext(Range(1, 4), Range(4, 7), Range(7, 10)),
			OnErrorResumeNext(Throw(errTest), Range(4, 7), Range(7, 10)),
			OnErrorResumeNext(Range(1, 4), Throw(errTest), Range(7, 10)),
			OnErrorResumeNext(Range(1, 4), Range(4, 7), Throw(errTest)),
		},
		[][]interface{}{
			{1, 2, 3, 4, 5, 6, 7, 8, 9, Complete},
			{4, 5, 6, 7, 8, 9, Complete},
			{1, 2, 3, 7, 8, 9, Complete},
			{1, 2, 3, 4, 5, 6, Complete},
		},
	)
}
