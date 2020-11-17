package rxtest

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
)

var (
	Completed = errors.New("completed")
	ErrTest   = errors.New("test")
)

func Step(n int) time.Duration {
	return 60 * time.Millisecond * time.Duration(n)
}

func AddLatencyToValues(initialDelay, period int) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		initialDelay, period := Step(initialDelay), Step(period)

		return rx.Zip(
			source,
			rx.Concat(rx.Timer(initialDelay), rx.Ticker(period)),
		).Pipe(
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

		return rx.Zip(
			source.Pipe(operators.Materialize()),
			rx.Concat(rx.Timer(initialDelay), rx.Ticker(period)),
		).Pipe(
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
		return rx.Timer(Step(n)).Pipe(operators.ConcatMapTo(source))
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

		_ = source.BlockingSubscribe(
			context.Background(),
			func(n rx.Notification) {
				if len(output) == 0 {
					t.Fail()
					switch {
					case n.HasValue:
						t.Logf("want nothing, but got %v", n.Value)
					case n.HasError:
						t.Logf("want nothing, but got %v", n.Error)
					default:
						t.Log("want nothing, but got completed")
					}
					return
				}

				wanted := output[0]
				output = output[1:]

				switch {
				case n.HasValue:
					if wanted != n.Value {
						t.Fail()
						t.Logf("want %v, but got %v", wanted, n.Value)
					} else {
						t.Logf("want %v", wanted)
					}
				case n.HasError:
					if wanted != n.Error {
						t.Fail()
						t.Logf("want %v, but got %v", wanted, n.Error)
					} else {
						t.Logf("want %v", wanted)
					}
				default:
					if wanted != Completed {
						t.Fail()
						t.Logf("want %v, but got completed", wanted)
					} else {
						t.Log("want completed")
					}
				}
			},
		)

		if len(output) > 0 {
			t.Fail()

			for _, wanted := range output {
				t.Logf("want %v, but got nothing", wanted)
			}
		}
	}
}
