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

	// To Subscribe to an observable, you call its Subscribe method with
	// a context and an observer as parameters, and you get another context
	// and a cancel function in return.
	ctx, cancel := obs.Subscribe(context.Background(), rx.Noop)

	// The returned context will be cancelled when the subscription completes.
	// Since this example has no goroutines involved, it must have already done.
	<-ctx.Done() // Wait for it done, though it's not necessary.

	// Subscriptions to observables are cancellable, but since this example has
	// already finished its work, there is nothing left to cancel.
	cancel() // This is a noop in this example.

	// To check if a subscription to an observable has really completed without
	// an error, we can call the Err method of the returned context to get
	// an error and check if it equals to rx.Completed.
	switch err := ctx.Err(); err {
	case nil:
		fmt.Println("Status: WIP")
	case rx.Completed:
		fmt.Println("Status: Completed")
	case context.Canceled:
		fmt.Println("Status: Canceled")
	default:
		fmt.Printf("Status: Error(%v)\n", err)
	}

	// Output:
	// 2
	// 6
	// 10
	// 14
	// 18
	// Completed
	// Status: Completed
}
