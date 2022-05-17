package rx

type observables[T any] []Observable[T]

func identity[T any](v T) T { return v }
