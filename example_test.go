package rx_test

import (
	"context"
	"fmt"

	"github.com/b97tsk/rx"
)

func Example() {
	// Create an Observable...
	obs := rx.Range(1, 10)

	// ...and apply some Operators.
	obs = rx.Pipe3(
		obs,
		rx.Filter(
			func(v int) bool {
				return v%2 == 1
			},
		),
		rx.Map(
			func(v int) int {
				return v * 2
			},
		),
		rx.Do(
			func(n rx.Notification[int]) {
				switch {
				case n.HasValue:
					fmt.Println(n.Value)
				case n.HasError:
					fmt.Println(n.Error)
				default:
					fmt.Println("Completed")
				}
			},
		),
	)

	// To Subscribe to an Observable, you call its Subscribe method which
	// takes a context.Context and an Observer as arguments.
	obs.Subscribe(context.TODO(), rx.Noop[int])

	// Since this example has no other goroutines involved, it must have
	// already done. You could also use BlockingSubscribe method instead.
	// It blocks until done.

	// Output:
	// 2
	// 6
	// 10
	// 14
	// 18
	// Completed
}
