package rx

import (
	"context"
)

type distinctOperator struct {
	source      Operator
	keySelector func(interface{}) interface{}
}

func (op distinctOperator) ApplyOptions(options []Option) Operator {
	for _, opt := range options {
		switch t := opt.(type) {
		case keySelectorOption:
			op.keySelector = t.Value
		default:
			panic(ErrUnsupportedOption)
		}
	}
	return op
}

func (op distinctOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	keys := make(map[interface{}]struct{})
	return op.source.Call(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			key := op.keySelector(t.Value)
			if _, exists := keys[key]; exists {
				break
			}
			keys[key] = struct{}{}
			ob.Next(t.Value)
		case t.HasError:
			ob.Error(t.Value.(error))
		default:
			ob.Complete()
		}
	})
}

// Distinct creates an Observable that emits all items emitted by the source
// Observable that are distinct by comparison from previous items.
//
// If a keySelector function is provided, then it will project each value from
// the source Observable into a new value that it will check for equality with
// previously projected values. If a keySelector function is not provided, it
// will use each value from the source Observable directly with an equality
// check against previous values.
func (o Observable) Distinct() Observable {
	op := distinctOperator{
		source:      o.Op,
		keySelector: defaultKeySelector,
	}
	return Observable{op}
}
