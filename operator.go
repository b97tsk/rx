package rx

import (
	"context"
)

// An Operator is the implementation logic of an OperatorFunc.
type Operator func(context.Context, Observer, Observable) (context.Context, context.CancelFunc)

// An OperatorFunc is an operation on an Observable. When called, they do not
// change the existing Observable instance. Instead, they return a new
// Observable, whose subscription logic is based on the first Observable.
type OperatorFunc func(Observable) Observable

// Pipe stitches Operators together into a chain.
func Pipe(operations ...OperatorFunc) OperatorFunc {
	return func(source Observable) Observable {
		return source.Pipe(operations...)
	}
}
