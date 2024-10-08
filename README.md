# rx

[![Go Reference](https://pkg.go.dev/badge/github.com/b97tsk/rx.svg)](https://pkg.go.dev/github.com/b97tsk/rx)
[![Build Status](https://github.com/b97tsk/rx/actions/workflows/build.yml/badge.svg?branch=generics)](https://github.com/b97tsk/rx/actions/workflows/build.yml)
[![Coverage Status](https://coveralls.io/repos/github/b97tsk/rx/badge.svg?branch=generics)](https://coveralls.io/github/b97tsk/rx?branch=generics)

A reactive programming library for Go, inspired by https://reactivex.io/ (mostly RxJS).

The code is in [another branch](https://github.com/b97tsk/rx/tree/generics).

Future updates are not guaranteed.

## Examples

[example_test.go](https://github.com/b97tsk/rx/blob/generics/example_test.go)

```go
// Create an Observable...
ob := rx.Range(1, 10)

// ...and apply some Operators.
ob = rx.Pipe3(
	ob,
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
			case rx.KindComplete:
				fmt.Println("Complete")
			case rx.KindError:
				fmt.Println("Error:", n.Error)
			case rx.KindStop:
				fmt.Println("Stop:", n.Error)
			}
		},
	),
)

// To Subscribe to an Observable, you call its Subscribe method, which takes
// a Context and an Observer as arguments.
ob.Subscribe(rx.NewContext(context.TODO()), rx.Noop[int])

// Since this example has no other goroutines involved, it must have already done.
// You can also use BlockingSubscribe method instead. It blocks until done.

// Output:
// 2
// 6
// 10
// 14
// 18
// Complete
```

```go
ctx, cancel := rx.NewContext(context.TODO()).WithTimeout(700 * time.Millisecond)
defer cancel()

ob := rx.Pipe4(
	rx.Concat(
		rx.Timer(50*time.Millisecond),
		rx.Ticker(100*time.Millisecond),
	),
	rx.Scan(-1, func(i int, _ time.Time) int { return i + 1 }), // 0, 1, 2, 3, ...
	rx.Scan(0, func(v1, v2 int) int { return v1 + v2 }),        // 0, 1, 3, 6, ...
	rx.Skip[int](4),
	rx.DoOnNext(func(v int) { fmt.Println(v) }), // 10, 15, 21, ...
)

switch n := ob.BlockingSubscribe(ctx, rx.Noop[int]); n.Kind {
case rx.KindComplete:
	fmt.Println("Complete")
case rx.KindError:
	fmt.Println("Error:", n.Error)
case rx.KindStop:
	fmt.Println("Stop:", n.Error)
}

// Output:
// 10
// 15
// 21
// Stop: context deadline exceeded
```

```go
ctx := rx.NewContext(context.TODO()).WithNewWaitGroup()

ctx.Go(func() {
	for n := 1; n <= 3; n++ {
		rx.Pipe2(
			rx.Timer(50*time.Millisecond*time.Duration(n)),
			rx.MapTo[time.Time](n),
			rx.DoOnNext(func(v int) { fmt.Println(v) }),
		).Subscribe(ctx, rx.Noop[int])
	}
})

ctx.Wait()

// Output:
// 1
// 2
// 3
```

```go
// Multicasts and Unicasts are special Observables that developers can decide
// what values they produce or when they complete, later after they are subscribed.
// Multicasts can be subscribed multiple times, whereas Unicasts can only be
// successfully subscribed once.
// Both Multicasts and Unicasts are safe for concurrent use.
// Here is an example demonstrates how to use a Unicast.

ctx := rx.NewContext(context.TODO()).WithNewWaitGroup()

u := rx.Unicast[time.Time]()

ctx.Go(func() {
	rx.Pipe2(
		rx.Concat(
			u.Observable,
			rx.Timer(500*time.Millisecond),
		),
		rx.MapTo[time.Time](42),
		rx.Do(
			func(n rx.Notification[int]) {
				switch n.Kind {
				case rx.KindNext:
					fmt.Println(n.Value)
				case rx.KindComplete:
					fmt.Println("Complete")
				case rx.KindError:
					fmt.Println("Error:", n.Error)
				case rx.KindStop:
					fmt.Println("Stop:", n.Error)
				}
			},
		),
	).Subscribe(ctx, rx.Noop[int])
})

ctx.Go(func() {
	u.Complete() // Start timer.
})

ctx.Wait()

// This example works properly, no matter which ctx.Go(...) runs first.

// Output:
// 42
// Complete
```
