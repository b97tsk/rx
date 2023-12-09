package rx_test

import (
	"context"
	"fmt"
	"time"

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
				switch n.Kind {
				case rx.KindNext:
					fmt.Println(n.Value)
				case rx.KindError:
					fmt.Println(n.Error)
				case rx.KindComplete:
					fmt.Println("Complete")
				}
			},
		),
	)

	// To Subscribe to an Observable, you call its Subscribe method, which takes
	// a context.Context and an Observer as arguments.
	obs.Subscribe(context.TODO(), rx.Noop[int])

	// Since this example has no other goroutines involved, it must have already done.
	// You can also use BlockingSubscribe method instead. It blocks until done.

	// Output:
	// 2
	// 6
	// 10
	// 14
	// 18
	// Complete
}

func Example_blocking() {
	ctx, cancel := context.WithTimeout(context.TODO(), 700*time.Millisecond)
	defer cancel()

	obs := rx.Pipe4(
		rx.Concat(
			rx.Timer(50*time.Millisecond),
			rx.Ticker(100*time.Millisecond),
		),
		rx.Scan(-1, func(i int, _ time.Time) int { return i + 1 }), // 0, 1, 2, 3, ...
		rx.Scan(0, func(v1, v2 int) int { return v1 + v2 }),        // 0, 1, 3, 6, ...
		rx.Skip[int](4),
		rx.OnNext(func(v int) { fmt.Println(v) }), // 10, 15, 21, ...
	)

	err := obs.BlockingSubscribe(ctx, rx.Noop[int])
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// 10
	// 15
	// 21
	// context deadline exceeded
}

func Example_waitGroup() {
	var wg rx.WaitGroup

	ctx := rx.WithWaitGroup(context.TODO(), &wg)

	wg.Go(func() {
		for n := 1; n < 4; n++ {
			rx.Pipe2(
				rx.Timer(50*time.Millisecond*time.Duration(n)),
				rx.MapTo[time.Time](n),
				rx.OnNext(func(v int) { fmt.Println(v) }),
			).Subscribe(ctx, rx.Noop[int])
		}
	})

	wg.Wait()

	// Output:
	// 1
	// 2
	// 3
}
