package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Congest(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Range(1, 9).Pipe(
				addLatencyToValue(1, 1),
				operators.Congest(3),
				addLatencyToValue(3, 4),
			),
		},
		1, 2, 3, 4, 5, 6, 7, 8, Complete,
	)
}
