package rx_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	. "github.com/b97tsk/rx"
)

var (
	xComplete = errors.New("complete")
	xErrTest  = errors.New("test")

	toString = operators.Map(
		func(val interface{}, idx int) interface{} {
			return fmt.Sprint(val)
		},
	)

	operators Operators
)

func step(n int) time.Duration {
	return 60 * time.Millisecond * time.Duration(n)
}

func addLatencyToValue(initialDelay, period int) OperatorFunc {
	return func(source Observable) Observable {
		initialDelay, period := step(initialDelay), step(period)
		return Zip(source, Timer(initialDelay, period)).Pipe(
			operators.Map(
				func(val interface{}, idx int) interface{} {
					return val.([]interface{})[0]
				},
			),
		)
	}
}

func addLatencyToNotification(initialDelay, period int) OperatorFunc {
	return func(source Observable) Observable {
		initialDelay, period := step(initialDelay), step(period)
		return Zip(source.Pipe(operators.Materialize()), Timer(initialDelay, period)).Pipe(
			operators.Map(
				func(val interface{}, idx int) interface{} {
					return val.([]interface{})[0]
				},
			),
			operators.Dematerialize(),
		)
	}
}

func delaySubscription(n int) OperatorFunc {
	return func(source Observable) Observable {
		return source.Pipe(operators.SubscribeOn(step(n)))
	}
}

func subscribe(tt *testing.T, observables []Observable, output ...interface{}) {
	for _, source := range observables {
		source.BlockingSubscribe(
			context.Background(),
			func(t Notification) {
				if len(output) == 0 {
					if t.HasValue || t.HasError {
						tt.Logf("expect Nothing, but got %v", t.Value)
					} else {
						tt.Log("expect Nothing, but got Complete")
					}
					tt.Fail()
					return
				}

				expected := output[0]
				output = output[1:]

				if t.HasValue || t.HasError {
					if expected != t.Value {
						tt.Logf("expect %v, but got %v", expected, t.Value)
						tt.Fail()
					} else {
						tt.Logf("expect %v", expected)
					}
				} else {
					if expected != xComplete {
						tt.Logf("expect %v, but got Complete", expected)
						tt.Fail()
					} else {
						tt.Log("expect Complete")
					}
				}
			},
		)
	}
	if len(output) > 0 {
		for _, expected := range output {
			tt.Logf("expect %v, but got Nothing", expected)
		}
		tt.Fail()
	}
}
