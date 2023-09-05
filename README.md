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
```

```go
ctx, cancel := context.WithTimeout(context.Background(), 700*time.Millisecond)
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
```
