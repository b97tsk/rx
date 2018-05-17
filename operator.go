package rx

import (
	"context"
)

// A Operator is a method on the Observable type. When called, they do not
// change the existing Observable instance. Instead, they return a new
// Observable, whose subscription logic is based on the first Observable.
type Operator func(context.Context, Observer, Observable) (context.Context, context.CancelFunc)
