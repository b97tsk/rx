package rx

// Compose2 creates an Operator that applies 2 Operators to the input
// Observable when called.
func Compose2[A, B, C any](
	op1 Operator[A, B],
	op2 Operator[B, C],
) Operator[A, C] {
	return NewOperator(
		func(obs Observable[A]) Observable[C] {
			return Pipe2(obs, op1, op2)
		},
	)
}

// Compose3 creates an Operator that applies 3 Operators to the input
// Observable when called.
func Compose3[A, B, C, D any](
	op1 Operator[A, B],
	op2 Operator[B, C],
	op3 Operator[C, D],
) Operator[A, D] {
	return NewOperator(
		func(obs Observable[A]) Observable[D] {
			return Pipe3(obs, op1, op2, op3)
		},
	)
}

// Compose4 creates an Operator that applies 4 Operators to the input
// Observable when called.
func Compose4[A, B, C, D, E any](
	op1 Operator[A, B],
	op2 Operator[B, C],
	op3 Operator[C, D],
	op4 Operator[D, E],
) Operator[A, E] {
	return NewOperator(
		func(obs Observable[A]) Observable[E] {
			return Pipe4(obs, op1, op2, op3, op4)
		},
	)
}

// Compose5 creates an Operator that applies 5 Operators to the input
// Observable when called.
func Compose5[A, B, C, D, E, F any](
	op1 Operator[A, B],
	op2 Operator[B, C],
	op3 Operator[C, D],
	op4 Operator[D, E],
	op5 Operator[E, F],
) Operator[A, F] {
	return NewOperator(
		func(obs Observable[A]) Observable[F] {
			return Pipe5(obs, op1, op2, op3, op4, op5)
		},
	)
}

// Compose6 creates an Operator that applies 6 Operators to the input
// Observable when called.
func Compose6[A, B, C, D, E, F, G any](
	op1 Operator[A, B],
	op2 Operator[B, C],
	op3 Operator[C, D],
	op4 Operator[D, E],
	op5 Operator[E, F],
	op6 Operator[F, G],
) Operator[A, G] {
	return NewOperator(
		func(obs Observable[A]) Observable[G] {
			return Pipe6(obs, op1, op2, op3, op4, op5, op6)
		},
	)
}

// Compose7 creates an Operator that applies 7 Operators to the input
// Observable when called.
func Compose7[A, B, C, D, E, F, G, H any](
	op1 Operator[A, B],
	op2 Operator[B, C],
	op3 Operator[C, D],
	op4 Operator[D, E],
	op5 Operator[E, F],
	op6 Operator[F, G],
	op7 Operator[G, H],
) Operator[A, H] {
	return NewOperator(
		func(obs Observable[A]) Observable[H] {
			return Pipe7(obs, op1, op2, op3, op4, op5, op6, op7)
		},
	)
}

// Compose8 creates an Operator that applies 8 Operators to the input
// Observable when called.
func Compose8[A, B, C, D, E, F, G, H, I any](
	op1 Operator[A, B],
	op2 Operator[B, C],
	op3 Operator[C, D],
	op4 Operator[D, E],
	op5 Operator[E, F],
	op6 Operator[F, G],
	op7 Operator[G, H],
	op8 Operator[H, I],
) Operator[A, I] {
	return NewOperator(
		func(obs Observable[A]) Observable[I] {
			return Pipe8(obs, op1, op2, op3, op4, op5, op6, op7, op8)
		},
	)
}

// Compose9 creates an Operator that applies 9 Operators to the input
// Observable when called.
func Compose9[A, B, C, D, E, F, G, H, I, J any](
	op1 Operator[A, B],
	op2 Operator[B, C],
	op3 Operator[C, D],
	op4 Operator[D, E],
	op5 Operator[E, F],
	op6 Operator[F, G],
	op7 Operator[G, H],
	op8 Operator[H, I],
	op9 Operator[I, J],
) Operator[A, J] {
	return NewOperator(
		func(obs Observable[A]) Observable[J] {
			return Pipe9(obs, op1, op2, op3, op4, op5, op6, op7, op8, op9)
		},
	)
}
