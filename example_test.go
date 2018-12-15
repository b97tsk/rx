package rx_test

import (
	"context"
	"fmt"

	"github.com/b97tsk/rx"
)

func Example() {
	var operators rx.Operators

	rx.Range(1, 10).Pipe(
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
					fmt.Println(t.Value)
				default:
					fmt.Println("Complete")
				}
			},
		),
	).Subscribe(context.Background(), rx.NopObserver)

	// Output:
	// 2
	// 6
	// 10
	// 14
	// 18
	// Complete
}
