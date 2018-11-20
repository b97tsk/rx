package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOnErrorResumeNext(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			OnErrorResumeNext(Range(1, 4), Range(4, 7), Range(7, 10)),
			OnErrorResumeNext(Throw(xErrTest), Range(4, 7), Range(7, 10)),
			OnErrorResumeNext(Range(1, 4), Throw(xErrTest), Range(7, 10)),
			OnErrorResumeNext(Range(1, 4), Range(4, 7), Throw(xErrTest)),
		},
		1, 2, 3, 4, 5, 6, 7, 8, 9, xComplete,
		4, 5, 6, 7, 8, 9, xComplete,
		1, 2, 3, 7, 8, 9, xComplete,
		1, 2, 3, 4, 5, 6, xComplete,
	)
}
