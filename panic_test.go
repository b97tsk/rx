package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
)

func TestPanic(t *testing.T) {
	t.Parallel()

	proj := func(v any) any { return v }

	shouldPanic(t, func() { _ = rx.Audit[any, any](nil) }, "Audit with durationSelector == nil")

	shouldPanic(t, func() { _ = rx.BufferCount[any](0) }, "BufferCount with bufferSize == 0")
	shouldPanic(t, func() {
		_ = rx.BufferCount[any](2).WithStartBufferEvery(0)
	}, "BufferCount with StartBufferEvery == 0")

	shouldPanic(t, func() { _ = rx.Catch[any](nil) }, "Catch with selector == nil")

	shouldPanic(t, func() { _ = rx.Channelize[any](nil) }, "Channelize with join == nil")

	shouldPanic(t, func() { _ = rx.CombineLatest2[any, any, any](nil, nil, nil) }, "CombineLatest2 with proj == nil")
	shouldPanic(t, func() {
		_ = rx.CombineLatest3[any, any, any, any](nil, nil, nil, nil)
	}, "CombineLatest3 with proj == nil")
	shouldPanic(t, func() {
		_ = rx.CombineLatest4[any, any, any, any, any](nil, nil, nil, nil, nil)
	}, "CombineLatest4 with proj == nil")
	shouldPanic(t, func() {
		_ = rx.CombineLatest5[any, any, any, any, any, any](nil, nil, nil, nil, nil, nil)
	}, "CombineLatest5 with proj == nil")
	shouldPanic(t, func() {
		_ = rx.CombineLatest6[any, any, any, any, any, any, any](nil, nil, nil, nil, nil, nil, nil)
	}, "CombineLatest6 with proj == nil")
	shouldPanic(t, func() {
		_ = rx.CombineLatest7[any, any, any, any, any, any, any, any](nil, nil, nil, nil, nil, nil, nil, nil)
	}, "CombineLatest7 with proj == nil")
	shouldPanic(t, func() {
		_ = rx.CombineLatest8[any, any, any, any, any, any, any, any, any](nil, nil, nil, nil, nil, nil, nil, nil, nil)
	}, "CombineLatest8 with proj == nil")
	shouldPanic(t, func() {
		f := rx.CombineLatest9[any, any, any, any, any, any, any, any, any, any]
		_ = f(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	}, "CombineLatest9 with proj == nil")

	shouldPanic(t, func() { _ = rx.Compact[any](nil) }, "Compact with eq == nil")
	shouldPanic(t, func() { _ = rx.CompactComparableKey[any, string](nil) }, "CompactComparableKey with proj == nil")
	shouldPanic(t, func() { _ = rx.CompactKey[any, any](nil, nil) }, "CompactKey with proj == nil")
	shouldPanic(t, func() { _ = rx.CompactKey(proj, nil) }, "CompactKey with eq == nil")

	shouldPanic(t, func() { _ = rx.ConcatMap[any, any](nil) }, "ConcatMap with proj == nil")

	shouldPanic(t, func() { _ = rx.CongestBlock[any](0) }, "CongestBlock with capacity == 0")
	shouldPanic(t, func() { _ = rx.CongestDropLatest[any](0) }, "CongestDropLatest with capacity == 0")
	shouldPanic(t, func() { _ = rx.CongestDropOldest[any](0) }, "CongestDropOldest with capacity == 0")
	shouldPanic(t, func() { _ = rx.CongestError[any](0) }, "CongestError with capacity == 0")

	shouldPanic(t, func() { _ = rx.Connect[any, any](nil) }, "Connect with selector == nil")
	shouldPanic(t, func() {
		selector := func(source rx.Observable[any]) rx.Observable[any] { return source }
		_ = rx.Connect(selector).WithConnector(nil)
	}, "Connect with Connector == nil")

	shouldPanic(t, func() { _ = rx.Contains[any](nil) }, "Contains with cond == nil")

	shouldPanic(t, func() { _ = rx.Debounce[any, any](nil) }, "Debounce with durationSelector == nil")

	shouldPanic(t, func() { _ = rx.Distinct[any, string](nil) }, "Distinct with proj == nil")

	shouldPanic(t, func() { _ = rx.Do[any](nil) }, "Do with tap == nil")
	shouldPanic(t, func() { _ = rx.OnNext[any](nil) }, "OnNext with f == nil")
	shouldPanic(t, func() { _ = rx.OnComplete[any](nil) }, "OnComplete with f == nil")
	shouldPanic(t, func() { _ = rx.OnError[any](nil) }, "OnError with f == nil")
	shouldPanic(t, func() { _ = rx.OnLastNotification[any](nil) }, "OnLastNotification with f == nil")

	shouldPanic(t, func() { _ = rx.Every[any](nil) }, "Every with cond == nil")

	shouldPanic(t, func() { _ = rx.ExhaustMap[any, any](nil) }, "ExhaustMap with proj == nil")

	shouldPanic(t, func() { _ = rx.Filter[any](nil) }, "Filter with cond == nil")
	shouldPanic(t, func() { _ = rx.FilterOut[any](nil) }, "FilterOut with cond == nil")
	shouldPanic(t, func() { _ = rx.FilterMap[any, any](nil) }, "FilterMap with cond == nil")

	shouldPanic(t, func() { _ = rx.Find[any](nil) }, "Find with cond == nil")

	shouldPanic(t, func() { _ = rx.Flat[rx.Observable[any]](nil) }, "Flat with f == nil")

	shouldPanic(t, func() { _ = rx.GroupBy[any, string](nil, nil) }, "GroupBy with keySelector == nil")
	shouldPanic(t, func() { _ = rx.GroupBy(func(any) string { return "" }, nil) }, "GroupBy with groupFactory == nil")

	shouldPanic(t, func() { _ = rx.Map[any, any](nil) }, "Map with proj == nil")

	shouldPanic(t, func() { _ = rx.MergeMap[any, any](nil) }, "MergeMap with proj == nil")
	shouldPanic(t, func() {
		_ = rx.MergeMapTo[any](rx.Just(42)).WithConcurrency(0)
	}, "MergeMapTo with Concurrency == 0")

	shouldPanic(t, func() { _ = rx.Reduce[any, any](nil, nil) }, "Reduce with accumulator == nil")
	shouldPanic(t, func() { _ = rx.Scan[any, any](nil, nil) }, "Scan with accumulator == nil")

	shouldPanic(t, func() { _ = rx.Share[any]().WithConnector(nil) }, "Share with Connector == nil")

	shouldPanic(t, func() { _ = rx.SkipWhile[any](nil) }, "SkipWhile with cond == nil")

	shouldPanic(t, func() { _ = rx.SwitchMap[any, any](nil) }, "SwitchMap with proj == nil")

	shouldPanic(t, func() { _ = rx.TakeWhile[any](nil) }, "TakeWhile with cond == nil")

	shouldPanic(t, func() { _ = rx.Throttle[any, any](nil) }, "Throttle with durationSelector == nil")

	shouldPanic(t, func() { _ = rx.Ticker(0) }, "Ticker with d == 0")

	shouldPanic(t, func() { _ = rx.WithLatestFrom1[any, any, any](nil, nil) }, "WithLatestFrom1 with proj == nil")
	shouldPanic(t, func() {
		_ = rx.WithLatestFrom2[any, any, any, any](nil, nil, nil)
	}, "WithLatestFrom2 with proj == nil")
	shouldPanic(t, func() {
		_ = rx.WithLatestFrom3[any, any, any, any, any](nil, nil, nil, nil)
	}, "WithLatestFrom3 with proj == nil")
	shouldPanic(t, func() {
		_ = rx.WithLatestFrom4[any, any, any, any, any, any](nil, nil, nil, nil, nil)
	}, "WithLatestFrom4 with proj == nil")
	shouldPanic(t, func() {
		_ = rx.WithLatestFrom5[any, any, any, any, any, any, any](nil, nil, nil, nil, nil, nil)
	}, "WithLatestFrom5 with proj == nil")
	shouldPanic(t, func() {
		_ = rx.WithLatestFrom6[any, any, any, any, any, any, any, any](nil, nil, nil, nil, nil, nil, nil)
	}, "WithLatestFrom6 with proj == nil")
	shouldPanic(t, func() {
		_ = rx.WithLatestFrom7[any, any, any, any, any, any, any, any, any](nil, nil, nil, nil, nil, nil, nil, nil)
	}, "WithLatestFrom7 with proj == nil")
	shouldPanic(t, func() {
		_ = rx.WithLatestFrom8[any, any, any, any, any, any, any, any, any, any](nil, nil, nil, nil, nil, nil, nil, nil, nil)
	}, "WithLatestFrom8 with proj == nil")

	shouldPanic(t, func() { _ = rx.Zip2[any, any, any](nil, nil, nil) }, "Zip2 with proj == nil")
	shouldPanic(t, func() { _ = rx.Zip3[any, any, any, any](nil, nil, nil, nil) }, "Zip3 with proj == nil")
	shouldPanic(t, func() { _ = rx.Zip4[any, any, any, any, any](nil, nil, nil, nil, nil) }, "Zip4 with proj == nil")
	shouldPanic(t, func() {
		_ = rx.Zip5[any, any, any, any, any, any](nil, nil, nil, nil, nil, nil)
	}, "Zip5 with proj == nil")
	shouldPanic(t, func() {
		_ = rx.Zip6[any, any, any, any, any, any, any](nil, nil, nil, nil, nil, nil, nil)
	}, "Zip6 with proj == nil")
	shouldPanic(t, func() {
		_ = rx.Zip7[any, any, any, any, any, any, any, any](nil, nil, nil, nil, nil, nil, nil, nil)
	}, "Zip7 with proj == nil")
	shouldPanic(t, func() {
		_ = rx.Zip8[any, any, any, any, any, any, any, any, any](nil, nil, nil, nil, nil, nil, nil, nil, nil)
	}, "Zip8 with proj == nil")
	shouldPanic(t, func() {
		_ = rx.Zip9[any, any, any, any, any, any, any, any, any, any](nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	}, "Zip9 with proj == nil")
}

func shouldPanic(t *testing.T, f func(), name string) {
	t.Helper()

	defer func() {
		if recover() == nil {
			t.Log(name, "did not panic.")
			t.Fail()
		}
	}()

	f()
}
