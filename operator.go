package rx

// An Operator is an operation on an [Observable]. When applied, they do not
// change the existing [Observable] value. Instead, they return a new one,
// whose subscription logic is based on the first [Observable].
type Operator[T, R any] interface {
	Apply(source Observable[T]) Observable[R]
}

// NewOperator creates an [Operator] from f.
func NewOperator[T, R any](f func(source Observable[T]) Observable[R]) Operator[T, R] {
	return operator[T, R]{f}
}

type operator[T, R any] struct {
	f func(source Observable[T]) Observable[R]
}

func (op operator[T, R]) Apply(source Observable[T]) Observable[R] {
	return op.f(source)
}
