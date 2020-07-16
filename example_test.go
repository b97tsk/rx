package rx_test

import (
	"context"
	"fmt"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
)

func Example() {
	// Create an observable and apply some operators.
	obs := rx.Range(1, 10).Pipe(
		operators.Filter(
			func(val interface{}, idx int) bool {
				return val.(int)%2 == 1
			},
		),
		operators.Map(
			func(val interface{}, idx int) interface{} {
				return val.(int) * 2
			},
		),
		operators.Do(
			func(t rx.Notification) {
				switch {
				case t.HasValue:
					fmt.Println(t.Value)
				case t.HasError:
					fmt.Println(t.Error)
				default:
					fmt.Println("Completed")
				}
			},
		),
	)

	// To Subscribe to an Observable, you call its Subscribe method which
	// takes a context.Context and an Observer as arguments.
	obs.Subscribe(context.Background(), rx.Noop)

	// Since this example has no goroutines involved, it must have already done.
	// You could also use BlockingSubscribe method instead. It blocks until done.

	// Output:
	// 2
	// 6
	// 10
	// 14
	// 18
	// Completed
}
