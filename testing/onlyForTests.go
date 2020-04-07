package testing

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
)

var ErrTest = errors.New("test")

func Step(n int) time.Duration {
	return 60 * time.Millisecond * time.Duration(n)
}

func AddLatencyToValues(initialDelay, period int) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		initialDelay, period := Step(initialDelay), Step(period)
		return rx.Zip(source, rx.Timer(initialDelay, period)).Pipe(
			operators.Map(
				func(val interface{}, idx int) interface{} {
					return val.([]interface{})[0]
				},
			),
		)
	}
}

func AddLatencyToNotifications(initialDelay, period int) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		initialDelay, period := Step(initialDelay), Step(period)
		return rx.Zip(source.Pipe(operators.Materialize()), rx.Timer(initialDelay, period)).Pipe(
			operators.Map(
				func(val interface{}, idx int) interface{} {
					return val.([]interface{})[0]
				},
			),
			operators.Dematerialize(),
		)
	}
}

func DelaySubscription(n int) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return rx.Concat(rx.Empty().Pipe(operators.Delay(Step(n))), source)
	}
}

func ToString() rx.Operator {
	return operators.Map(
		func(val interface{}, idx int) interface{} {
			return fmt.Sprint(val)
		},
	)
}

func Subscribe(t *testing.T, obs rx.Observable, output ...interface{}) {
	SubscribeN(t, []rx.Observable{obs}, [][]interface{}{output})
}

func SubscribeN(t *testing.T, observables []rx.Observable, outputs [][]interface{}) {
	if len(observables) != len(outputs) {
		panic("SubscribeN: len(observables) != len(outputs)")
	}
	for i, source := range observables {
		output := outputs[i]
		source.BlockingSubscribe(
			context.Background(),
			func(x rx.Notification) {
				if len(output) == 0 {
					switch {
					case x.HasValue:
						t.Logf("expect nothing, but got %v", x.Value)
					case x.HasError:
						t.Logf("expect nothing, but got %v", x.Error)
					default:
						t.Log("expect nothing, but got complete")
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
					if expected != rx.Complete {
						t.Logf("expect %v, but got complete", expected)
						t.Fail()
					} else {
						t.Log("expect complete")
					}
				}
			},
		)
		if len(output) > 0 {
			for _, expected := range output {
				t.Logf("expect %v, but got nothing", expected)
			}
			t.Fail()
		}
	}
}
