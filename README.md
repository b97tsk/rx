# rx

[![Go Reference](https://pkg.go.dev/badge/github.com/b97tsk/rx.svg)](https://pkg.go.dev/github.com/b97tsk/rx)
[![Build Status](https://github.com/b97tsk/rx/actions/workflows/build.yml/badge.svg?branch=generics)](https://github.com/b97tsk/rx/actions)
[![codecov](https://codecov.io/gh/b97tsk/rx/branch/generics/graph/badge.svg?token=N7LYL4U1U8)](https://codecov.io/gh/b97tsk/rx)

A reactive programming library for Go, inspired by https://reactivex.io/ (mostly RxJS).

The code is in [another branch](https://github.com/b97tsk/rx/tree/generics).

Future updates are not guaranteed.

## Example

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
				fmt.Println("Completed")
			}
		},
	),
)

// To Subscribe to an Observable, you call its Subscribe method which
// takes a context.Context and an Observer as arguments.
obs.Subscribe(context.TODO(), rx.Noop[int])
```
