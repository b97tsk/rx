package rx_test

import (
	"context"
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Multicast(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Interval(step(1)).Pipe(
				operators.Multicast(
					func(ctx context.Context) Subject {
						subject := NewSubject()
						return Subject{
							Zip(
								subject.Pipe(operators.Take(4)),
								subject.Pipe(operators.Skip(4), operators.Take(4)),
							).Pipe(
								operators.Map(
									func(val interface{}, idx int) interface{} {
										vals := val.([]interface{})
										return vals[0].(int) * vals[1].(int)
									},
								),
							),
							subject.Observer,
						}
					},
				),
			),
		},
		0, 5, 12, 21, xComplete,
	)
}
