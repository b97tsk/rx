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
	operators Operators

	toString = operators.Map(
		func(val interface{}, idx int) interface{} {
			return fmt.Sprint(val)
		},
	)

	errTest = errors.New("errTest")
)

func step(n int) time.Duration {
	return 60 * time.Millisecond * time.Duration(n)
}

func addLatencyToValue(initialDelay, period int) Operator {
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

func addLatencyToNotification(initialDelay, period int) Operator {
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

func delaySubscription(n int) Operator {
	return func(source Observable) Observable {
		return source.Pipe(operators.SubscribeOn(step(n)))
	}
}

func subscribe(t *testing.T, observables []Observable, output ...interface{}) {
	for _, source := range observables {
		source.BlockingSubscribe(
			context.Background(),
			func(x Notification) {
				if len(output) == 0 {
					switch {
					case x.HasValue:
						t.Logf("expect Nothing, but got %v", x.Value)
					case x.HasError:
						t.Logf("expect Nothing, but got %v", x.Error)
					default:
						t.Log("expect Nothing, but got Complete")
					}
					t.Fail()
					return
				}

				expected := output[0]
				output = output[1:]

				switch {
				case x.HasValue:
					if expected != x.Value {
						t.Logf("expect %v, but got %v", expected, x.Value)
						t.Fail()
					} else {
						t.Logf("expect %v", expected)
					}
				case x.HasError:
					if expected != x.Error {
						t.Logf("expect %v, but got %v", expected, x.Error)
						t.Fail()
					} else {
						t.Logf("expect %v", expected)
					}
				default:
					if expected != Complete {
						t.Logf("expect %v, but got Complete", expected)
						t.Fail()
					} else {
						t.Log("expect Complete")
					}
				}
			},
		)
	}
	if len(output) > 0 {
		for _, expected := range output {
			t.Logf("expect %v, but got Nothing", expected)
		}
		t.Fail()
	}
}
