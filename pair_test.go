package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestFromMap(t *testing.T) {
	t.Parallel()

	m := map[int]string{3: "C", 1: "A", 2: "B"}

	add := func(m map[int]string, p rx.Pair[int, string]) map[int]string {
		m[p.Key] = p.Value
		return m
	}

	NewTestSuite[string](t).Case(
		rx.Pipe2(
			rx.FromMap(m),
			rx.Reduce(make(map[int]string), add),
			ToString[map[int]string](),
		),
		"map[1:A 2:B 3:C]", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.FromMap(m),
			rx.Single[rx.Pair[int, string]](),
			rx.ValueOf[rx.Pair[int, string]](),
		),
		rx.ErrNotSingle,
	)
}

func TestKeyOf(t *testing.T) {
	t.Parallel()

	sum := func(v1, v2 int) int {
		return v1 + v2
	}

	m := map[int]string{3: "C", 1: "A", 2: "B"}

	NewTestSuite[int](t).Case(
		rx.Pipe2(
			rx.FromMap(m),
			rx.KeyOf[rx.Pair[int, string]](),
			rx.Reduce(0, sum),
		),
		6, ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Throw[rx.Pair[int, string]](ErrTest),
			rx.KeyOf[rx.Pair[int, string]](),
			rx.Reduce(0, sum),
		),
		ErrTest,
	)
}

func TestValueOf(t *testing.T) {
	t.Parallel()

	max := func(v1, v2 string) string {
		if v1 > v2 {
			return v1
		}

		return v2
	}

	m := map[int]string{3: "C", 1: "A", 2: "B"}

	NewTestSuite[string](t).Case(
		rx.Pipe2(
			rx.FromMap(m),
			rx.ValueOf[rx.Pair[int, string]](),
			rx.Reduce("", max),
		),
		"C", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Throw[rx.Pair[int, string]](ErrTest),
			rx.ValueOf[rx.Pair[int, string]](),
			rx.Reduce("", max),
		),
		ErrTest,
	)
}

func TestWithIndex(t *testing.T) {
	t.Parallel()

	add := func(m map[int]string, p rx.Pair[int, string]) map[int]string {
		m[p.Key] = p.Value
		return m
	}

	NewTestSuite[string](t).Case(
		rx.Pipe3(
			rx.Just("A", "B", "C"),
			rx.WithIndex[string](1),
			rx.Reduce(make(map[int]string), add),
			ToString[map[int]string](),
		),
		"map[1:A 2:B 3:C]", ErrComplete,
	).Case(
		rx.Pipe3(
			rx.Throw[string](ErrTest),
			rx.WithIndex[string](1),
			rx.Reduce(make(map[int]string), add),
			ToString[map[int]string](),
		),
		ErrTest,
	)
}
