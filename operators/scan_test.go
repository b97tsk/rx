package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestScan(t *testing.T) {
	max := func(acc, val interface{}, idx int) interface{} {
		if acc.(int) > val.(int) {
			return acc
		}

		return val
	}

	sum := func(acc, val interface{}, idx int) interface{} {
		return acc.(int) + val.(int)
	}

	NewTestSuite(t).Case(
		rx.Range(1, 7).Pipe(
			operators.Scan(max),
		),
		1, 2, 3, 4, 5, 6, Completed,
	).Case(
		rx.Just(42).Pipe(
			operators.Scan(max),
		),
		42, Completed,
	).Case(
		rx.Empty().Pipe(
			operators.Scan(max),
		),
		Completed,
	).Case(
		rx.Range(1, 7).Pipe(
			operators.Scan(sum),
		),
		1, 3, 6, 10, 15, 21, Completed,
	).Case(
		rx.Just(42).Pipe(
			operators.Scan(sum),
		),
		42, Completed,
	).Case(
		rx.Empty().Pipe(
			operators.Scan(sum),
		),
		Completed,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.Scan(sum),
		),
		ErrTest,
	).TestAll()

	panictest := func(f func(), msg string) {
		defer func() {
			if recover() == nil {
				t.Log(msg)
				t.FailNow()
			}
		}()
		f()
	}
	panictest(
		func() {
			operators.ScanConfigure{
				Accumulator: nil,
			}.Make()
		},
		"ScanConfigure with nil Accumulator didn't panic.",
	)
}
