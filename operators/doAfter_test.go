package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
)

func TestDoAfter(t *testing.T) {
	t.Run("DoAfter", func(t *testing.T) {
		testDo(
			t,
			func(f func()) rx.Operator {
				return operators.DoAfter(func(rx.Notification) { f() })
			},
			1, 3, 6, 9,
		)
	})
	t.Run("DoAfterNext", func(t *testing.T) {
		testDo(
			t,
			func(f func()) rx.Operator {
				return operators.DoAfterNext(func(interface{}) { f() })
			},
			0, 1, 3, 5,
		)
	})
	t.Run("DoAfterError", func(t *testing.T) {
		testDo(
			t,
			func(f func()) rx.Operator {
				return operators.DoAfterError(func(error) { f() })
			},
			0, 0, 0, 1,
		)
	})
	t.Run("DoAfterComplete", func(t *testing.T) {
		testDo(t, operators.DoAfterComplete, 1, 2, 3, 3)
	})
	t.Run("DoAfterErrorOrComplete", func(t *testing.T) {
		testDo(t, operators.DoAfterErrorOrComplete, 1, 2, 3, 4)
	})
}
