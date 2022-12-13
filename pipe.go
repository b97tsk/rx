package rx

// Pipe1 applies an Operator to an Observable and returns the resulting
// Observable.
func Pipe1[A, B any](
	obs Observable[A],
	op1 Operator[A, B],
) Observable[B] {
	return op1.Apply(obs)
}

// Pipe2 applies 2 Operators to an Observable and returns the resulting
// Observable.
func Pipe2[A, B, C any](
	obs Observable[A],
	op1 Operator[A, B],
	op2 Operator[B, C],
) Observable[C] {
	return op2.Apply(op1.Apply(obs))
}

// Pipe3 applies 3 Operators to an Observable and returns the resulting
// Observable.
func Pipe3[A, B, C, D any](
	obs Observable[A],
	op1 Operator[A, B],
	op2 Operator[B, C],
	op3 Operator[C, D],
) Observable[D] {
	return op3.Apply(op2.Apply(op1.Apply(obs)))
}

// Pipe4 applies 4 Operators to an Observable and returns the resulting
// Observable.
func Pipe4[A, B, C, D, E any](
	obs Observable[A],
	op1 Operator[A, B],
	op2 Operator[B, C],
	op3 Operator[C, D],
	op4 Operator[D, E],
) Observable[E] {
	return op4.Apply(op3.Apply(op2.Apply(op1.Apply(obs))))
}

// Pipe5 applies 5 Operators to an Observable and returns the resulting
// Observable.
func Pipe5[A, B, C, D, E, F any](
	obs Observable[A],
	op1 Operator[A, B],
	op2 Operator[B, C],
	op3 Operator[C, D],
	op4 Operator[D, E],
	op5 Operator[E, F],
) Observable[F] {
	return op5.Apply(op4.Apply(op3.Apply(op2.Apply(op1.Apply(obs)))))
}

// Pipe6 applies 6 Operators to an Observable and returns the resulting
// Observable.
func Pipe6[A, B, C, D, E, F, G any](
	obs Observable[A],
	op1 Operator[A, B],
	op2 Operator[B, C],
	op3 Operator[C, D],
	op4 Operator[D, E],
	op5 Operator[E, F],
	op6 Operator[F, G],
) Observable[G] {
	return op6.Apply(op5.Apply(op4.Apply(op3.Apply(op2.Apply(op1.Apply(obs))))))
}

// Pipe7 applies 7 Operators to an Observable and returns the resulting
// Observable.
func Pipe7[A, B, C, D, E, F, G, H any](
	obs Observable[A],
	op1 Operator[A, B],
	op2 Operator[B, C],
	op3 Operator[C, D],
	op4 Operator[D, E],
	op5 Operator[E, F],
	op6 Operator[F, G],
	op7 Operator[G, H],
) Observable[H] {
	return op7.Apply(op6.Apply(op5.Apply(op4.Apply(op3.Apply(op2.Apply(op1.Apply(obs)))))))
}

// Pipe8 applies 8 Operators to an Observable and returns the resulting
// Observable.
func Pipe8[A, B, C, D, E, F, G, H, I any](
	obs Observable[A],
	op1 Operator[A, B],
	op2 Operator[B, C],
	op3 Operator[C, D],
	op4 Operator[D, E],
	op5 Operator[E, F],
	op6 Operator[F, G],
	op7 Operator[G, H],
	op8 Operator[H, I],
) Observable[I] {
	return op8.Apply(op7.Apply(op6.Apply(op5.Apply(op4.Apply(op3.Apply(op2.Apply(op1.Apply(obs))))))))
}

// Pipe9 applies 9 Operators to an Observable and returns the resulting
// Observable.
func Pipe9[A, B, C, D, E, F, G, H, I, J any](
	obs Observable[A],
	op1 Operator[A, B],
	op2 Operator[B, C],
	op3 Operator[C, D],
	op4 Operator[D, E],
	op5 Operator[E, F],
	op6 Operator[F, G],
	op7 Operator[G, H],
	op8 Operator[H, I],
	op9 Operator[I, J],
) Observable[J] {
	return op9.Apply(op8.Apply(op7.Apply(op6.Apply(op5.Apply(op4.Apply(op3.Apply(op2.Apply(op1.Apply(obs)))))))))
}
