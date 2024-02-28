package rx

// Try0 runs f().
// If f does not return normally, runs oops().
func Try0(f, oops func()) {
	var ok bool

	defer func() {
		if !ok {
			oops()
		}
	}()

	f()

	ok = true
}

// Try1 runs f(v1).
// If f does not return normally, runs oops().
func Try1[T1 any](f func(v1 T1), v1 T1, oops func()) {
	var ok bool

	defer func() {
		if !ok {
			oops()
		}
	}()

	f(v1)

	ok = true
}

// Try2 runs f(v1, v2).
// If f does not return normally, runs oops().
func Try2[T1, T2 any](f func(v1 T1, v2 T2), v1 T1, v2 T2, oops func()) {
	var ok bool

	defer func() {
		if !ok {
			oops()
		}
	}()

	f(v1, v2)

	ok = true
}

// Try3 runs f(v1, v2, v3).
// If f does not return normally, runs oops().
func Try3[T1, T2, T3 any](f func(v1 T1, v2 T2, v3 T3), v1 T1, v2 T2, v3 T3, oops func()) {
	var ok bool

	defer func() {
		if !ok {
			oops()
		}
	}()

	f(v1, v2, v3)

	ok = true
}

// Try4 runs f(v1, v2, v3, v4).
// If f does not return normally, runs oops().
func Try4[T1, T2, T3, T4 any](f func(v1 T1, v2 T2, v3 T3, v4 T4), v1 T1, v2 T2, v3 T3, v4 T4, oops func()) {
	var ok bool

	defer func() {
		if !ok {
			oops()
		}
	}()

	f(v1, v2, v3, v4)

	ok = true
}

// Try5 runs f(v1, v2, v3, v4, v5).
// If f does not return normally, runs oops().
func Try5[T1, T2, T3, T4, T5 any](
	f func(v1 T1, v2 T2, v3 T3, v4 T4, v5 T5),
	v1 T1, v2 T2, v3 T3, v4 T4, v5 T5,
	oops func(),
) {
	var ok bool

	defer func() {
		if !ok {
			oops()
		}
	}()

	f(v1, v2, v3, v4, v5)

	ok = true
}

// Try6 runs f(v1, v2, v3, v4, v5, v6).
// If f does not return normally, runs oops().
func Try6[T1, T2, T3, T4, T5, T6 any](
	f func(v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6),
	v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6,
	oops func(),
) {
	var ok bool

	defer func() {
		if !ok {
			oops()
		}
	}()

	f(v1, v2, v3, v4, v5, v6)

	ok = true
}

// Try7 runs f(v1, v2, v3, v4, v5, v6, v7).
// If f does not return normally, runs oops().
func Try7[T1, T2, T3, T4, T5, T6, T7 any](
	f func(v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7),
	v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7,
	oops func(),
) {
	var ok bool

	defer func() {
		if !ok {
			oops()
		}
	}()

	f(v1, v2, v3, v4, v5, v6, v7)

	ok = true
}

// Try8 runs f(v1, v2, v3, v4, v5, v6, v7, v8).
// If f does not return normally, runs oops().
func Try8[T1, T2, T3, T4, T5, T6, T7, T8 any](
	f func(v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7, v8 T8),
	v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7, v8 T8,
	oops func(),
) {
	var ok bool

	defer func() {
		if !ok {
			oops()
		}
	}()

	f(v1, v2, v3, v4, v5, v6, v7, v8)

	ok = true
}

// Try9 runs f(v1, v2, v3, v4, v5, v6, v7, v8, v9).
// If f does not return normally, runs oops().
func Try9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](
	f func(v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7, v8 T8, v9 T9),
	v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7, v8 T8, v9 T9,
	oops func(),
) {
	var ok bool

	defer func() {
		if !ok {
			oops()
		}
	}()

	f(v1, v2, v3, v4, v5, v6, v7, v8, v9)

	ok = true
}

// Try01 returns f().
// If f does not return normally, runs oops().
func Try01[R any](f func() R, oops func()) R {
	var ok bool

	defer func() {
		if !ok {
			oops()
		}
	}()

	v := f()

	ok = true

	return v
}

// Try11 returns f(v1).
// If f does not return normally, runs oops().
func Try11[T1, R any](f func(v1 T1) R, v1 T1, oops func()) R {
	var ok bool

	defer func() {
		if !ok {
			oops()
		}
	}()

	v := f(v1)

	ok = true

	return v
}

// Try21 returns f(v1, v2).
// If f does not return normally, runs oops().
func Try21[T1, T2, R any](f func(v1 T1, v2 T2) R, v1 T1, v2 T2, oops func()) R {
	var ok bool

	defer func() {
		if !ok {
			oops()
		}
	}()

	v := f(v1, v2)

	ok = true

	return v
}

// Try31 returns f(v1, v2, v3).
// If f does not return normally, runs oops().
func Try31[T1, T2, T3, R any](f func(v1 T1, v2 T2, v3 T3) R, v1 T1, v2 T2, v3 T3, oops func()) R {
	var ok bool

	defer func() {
		if !ok {
			oops()
		}
	}()

	v := f(v1, v2, v3)

	ok = true

	return v
}

// Try41 returns f(v1, v2, v3, v4).
// If f does not return normally, runs oops().
func Try41[T1, T2, T3, T4, R any](f func(v1 T1, v2 T2, v3 T3, v4 T4) R, v1 T1, v2 T2, v3 T3, v4 T4, oops func()) R {
	var ok bool

	defer func() {
		if !ok {
			oops()
		}
	}()

	v := f(v1, v2, v3, v4)

	ok = true

	return v
}

// Try51 returns f(v1, v2, v3, v4, v5).
// If f does not return normally, runs oops().
func Try51[T1, T2, T3, T4, T5, R any](
	f func(v1 T1, v2 T2, v3 T3, v4 T4, v5 T5) R,
	v1 T1, v2 T2, v3 T3, v4 T4, v5 T5,
	oops func(),
) R {
	var ok bool

	defer func() {
		if !ok {
			oops()
		}
	}()

	v := f(v1, v2, v3, v4, v5)

	ok = true

	return v
}

// Try61 returns f(v1, v2, v3, v4, v5, v6).
// If f does not return normally, runs oops().
func Try61[T1, T2, T3, T4, T5, T6, R any](
	f func(v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6) R,
	v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6,
	oops func(),
) R {
	var ok bool

	defer func() {
		if !ok {
			oops()
		}
	}()

	v := f(v1, v2, v3, v4, v5, v6)

	ok = true

	return v
}

// Try71 returns f(v1, v2, v3, v4, v5, v6, v7).
// If f does not return normally, runs oops().
func Try71[T1, T2, T3, T4, T5, T6, T7, R any](
	f func(v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7) R,
	v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7,
	oops func(),
) R {
	var ok bool

	defer func() {
		if !ok {
			oops()
		}
	}()

	v := f(v1, v2, v3, v4, v5, v6, v7)

	ok = true

	return v
}

// Try81 returns f(v1, v2, v3, v4, v5, v6, v7, v8).
// If f does not return normally, runs oops().
func Try81[T1, T2, T3, T4, T5, T6, T7, T8, R any](
	f func(v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7, v8 T8) R,
	v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7, v8 T8,
	oops func(),
) R {
	var ok bool

	defer func() {
		if !ok {
			oops()
		}
	}()

	v := f(v1, v2, v3, v4, v5, v6, v7, v8)

	ok = true

	return v
}

// Try91 returns f(v1, v2, v3, v4, v5, v6, v7, v8, v9).
// If f does not return normally, runs oops().
func Try91[T1, T2, T3, T4, T5, T6, T7, T8, T9, R any](
	f func(v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7, v8 T8, v9 T9) R,
	v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7, v8 T8, v9 T9,
	oops func(),
) R {
	var ok bool

	defer func() {
		if !ok {
			oops()
		}
	}()

	v := f(v1, v2, v3, v4, v5, v6, v7, v8, v9)

	ok = true

	return v
}
