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

// MakeFunc casts an Operator into an OperatorFunc.
func MakeFunc(op Operator) OperatorFunc {
	return func(source Observable) Observable {
		return source.Lift(op)
	}
}
