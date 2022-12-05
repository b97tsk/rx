package rx

// An Operator is an operation on an Observable. When called, they do not
// change the existing Observable instance. Instead, they return a new
// Observable, whose subscription logic is based on the first Observable.
type Operator[T, R any] interface {
	Apply(source Observable[T]) Observable[R]
}

// NewOperator creates an Operator from f.
func NewOperator[T, R any](f func(source Observable[T]) Observable[R]) Operator[T, R] {
	return operatorFunc[T, R](f)
}

type operatorFunc[T, R any] func(source Observable[T]) Observable[R]

func (f operatorFunc[T, R]) Apply(source Observable[T]) Observable[R] {
	return f(source)
}
