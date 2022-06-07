package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
)

func TestPanic(t *testing.T) {
	t.Parallel()

	proj := func(v any) any { return v }

	shouldPanic(t, func() { _ = rx.Catch[any](nil) }, "Catch with selector == nil")
	shouldPanic(t, func() { _ = rx.OnErrorResumeWith[any](nil) }, "OnErrorResumeWith with obs == nil")

	shouldPanic(t, func() { _ = rx.Compact[any](nil) }, "Compact with eq == nil")
	shouldPanic(t, func() { _ = rx.CompactComparableKey[any, string](nil) }, "CompactComparableKey with proj == nil")
	shouldPanic(t, func() { _ = rx.CompactKey[any, any](nil, nil) }, "CompactKey with proj == nil")
	shouldPanic(t, func() { _ = rx.CompactKey(proj, nil) }, "CompactKey with eq == nil")

	shouldPanic(t, func() { _ = rx.ConcatMap[any, any](nil) }, "ConcatMap with proj == nil")
	shouldPanic(t, func() { _ = rx.ConcatMapTo[any, any](nil) }, "ConcatMapTo with inner == nil")

	shouldPanic(t, func() { _ = rx.Contains[any](nil) }, "Contains with cond == nil")

	shouldPanic(t, func() { _ = rx.Distinct[any, string](nil) }, "Distinct with proj == nil")

	shouldPanic(t, func() { _ = rx.Do[any](nil) }, "Do with tap == nil")
	shouldPanic(t, func() { _ = rx.DoOnNext[any](nil) }, "DoOnNext with f == nil")
	shouldPanic(t, func() { _ = rx.DoOnComplete[any](nil) }, "DoOnComplete with f == nil")
	shouldPanic(t, func() { _ = rx.DoOnError[any](nil) }, "DoOnError with f == nil")
	shouldPanic(t, func() { _ = rx.DoOnErrorOrComplete[any](nil) }, "DoOnErrorOrComplete with f == nil")

	shouldPanic(t, func() { _ = rx.Every[any](nil) }, "Every with cond == nil")

	shouldPanic(t, func() { _ = rx.ExhaustMap[any, any](nil) }, "ExhaustMap with proj == nil")
	shouldPanic(t, func() { _ = rx.ExhaustMapTo[any, any](nil) }, "ExhaustMapTo with inner == nil")

	shouldPanic(t, func() { _ = rx.Filter[any](nil) }, "Filter with cond == nil")
	shouldPanic(t, func() { _ = rx.FilterOut[any](nil) }, "FilterOut with cond == nil")
	shouldPanic(t, func() { _ = rx.FilterMap[any, any](nil) }, "FilterMap with cond == nil")

	shouldPanic(t, func() { _ = rx.Find[any](nil) }, "Find with cond == nil")

	shouldPanic(t, func() { _ = rx.Flat[rx.Observable[any]](nil) }, "Flat with f == nil")

	shouldPanic(t, func() { _ = rx.Map[any, any](nil) }, "Map with proj == nil")

	shouldPanic(t, func() { _ = rx.MergeMap[any, any](nil) }, "MergeMap with proj == nil")
	shouldPanic(t, func() { _ = rx.MergeMapTo[any, any](nil) }, "MergeMapTo with inner == nil")
	shouldPanic(t, func() {
		_ = rx.MergeMapTo[any](rx.Just(42)).WithConcurrency(0)
	}, "MergeMapTo with Concurrency == 0")

	shouldPanic(t, func() { _ = rx.Reduce[any, any](nil, nil) }, "Reduce with accumulator == nil")
	shouldPanic(t, func() { _ = rx.Scan[any, any](nil, nil) }, "Scan with accumulator == nil")

	shouldPanic(t, func() { _ = rx.SkipUntil[any, any](nil) }, "SkipUntil with notifier == nil")
	shouldPanic(t, func() { _ = rx.SkipWhile[any](nil) }, "SkipWhile with cond == nil")

	shouldPanic(t, func() { _ = rx.SwitchIfEmpty[any](nil) }, "SwitchIfEmpty with obs == nil")

	shouldPanic(t, func() { _ = rx.SwitchMap[any, any](nil) }, "SwitchMap with proj == nil")
	shouldPanic(t, func() { _ = rx.SwitchMapTo[any, any](nil) }, "SwitchMapTo with inner == nil")

	shouldPanic(t, func() { _ = rx.TakeUntil[any, any](nil) }, "TakeUntil with notifier == nil")
	shouldPanic(t, func() { _ = rx.TakeWhile[any](nil) }, "TakeWhile with cond == nil")

	shouldPanic(t, func() { _ = rx.Ticker(0) }, "Ticker with d == 0")

	panicTestZipLike(
		t,
		"CombineLatest",
		rx.CombineLatest2[any, any, any],
		rx.CombineLatest3[any, any, any, any],
		rx.CombineLatest4[any, any, any, any, any],
		rx.CombineLatest5[any, any, any, any, any, any],
		rx.CombineLatest6[any, any, any, any, any, any, any],
		rx.CombineLatest7[any, any, any, any, any, any, any, any],
		rx.CombineLatest8[any, any, any, any, any, any, any, any, any],
		rx.CombineLatest9[any, any, any, any, any, any, any, any, any, any],
	)

	panicTestZipLike(
		t,
		"Zip",
		rx.Zip2[any, any, any],
		rx.Zip3[any, any, any, any],
		rx.Zip4[any, any, any, any, any],
		rx.Zip5[any, any, any, any, any, any],
		rx.Zip6[any, any, any, any, any, any, any],
		rx.Zip7[any, any, any, any, any, any, any, any],
		rx.Zip8[any, any, any, any, any, any, any, any, any],
		rx.Zip9[any, any, any, any, any, any, any, any, any, any],
	)
}

func panicTestZipLike(
	t *testing.T,
	name string,
	f2 func(_, _ rx.Observable[any], _ func(any, any) any) rx.Observable[any],
	f3 func(_, _, _ rx.Observable[any], _ func(any, any, any) any) rx.Observable[any],
	f4 func(_, _, _, _ rx.Observable[any], _ func(any, any, any, any) any) rx.Observable[any],
	f5 func(_, _, _, _, _ rx.Observable[any], _ func(any, any, any, any, any) any) rx.Observable[any],
	f6 func(_, _, _, _, _, _ rx.Observable[any], _ func(any, any, any, any, any, any) any) rx.Observable[any],
	f7 func(_, _, _, _, _, _, _ rx.Observable[any], _ func(any, any, any, any, any, any, any) any) rx.Observable[any],
	f8 func(
		_, _, _, _, _, _, _, _ rx.Observable[any],
		_ func(any, any, any, any, any, any, any, any) any,
	) rx.Observable[any],
	f9 func(
		_, _, _, _, _, _, _, _, _ rx.Observable[any],
		_ func(any, any, any, any, any, any, any, any, any) any,
	) rx.Observable[any],
) {
	obs := rx.Empty[any]()

	shouldPanic(t, func() { _ = f2(nil, nil, nil) }, name+"2 with obs1 == nil")
	shouldPanic(t, func() { _ = f2(obs, nil, nil) }, name+"2 with obs2 == nil")
	shouldPanic(t, func() { _ = f2(obs, obs, nil) }, name+"2 with proj == nil")

	shouldPanic(t, func() { _ = f3(nil, nil, nil, nil) }, name+"3 with obs1 == nil")
	shouldPanic(t, func() { _ = f3(obs, nil, nil, nil) }, name+"3 with obs2 == nil")
	shouldPanic(t, func() { _ = f3(obs, obs, nil, nil) }, name+"3 with obs3 == nil")
	shouldPanic(t, func() { _ = f3(obs, obs, obs, nil) }, name+"3 with proj == nil")

	shouldPanic(t, func() { _ = f4(nil, nil, nil, nil, nil) }, name+"4 with obs1 == nil")
	shouldPanic(t, func() { _ = f4(obs, nil, nil, nil, nil) }, name+"4 with obs2 == nil")
	shouldPanic(t, func() { _ = f4(obs, obs, nil, nil, nil) }, name+"4 with obs3 == nil")
	shouldPanic(t, func() { _ = f4(obs, obs, obs, nil, nil) }, name+"4 with obs4 == nil")
	shouldPanic(t, func() { _ = f4(obs, obs, obs, obs, nil) }, name+"4 with proj == nil")

	shouldPanic(t, func() { _ = f5(nil, nil, nil, nil, nil, nil) }, name+"5 with obs1 == nil")
	shouldPanic(t, func() { _ = f5(obs, nil, nil, nil, nil, nil) }, name+"5 with obs2 == nil")
	shouldPanic(t, func() { _ = f5(obs, obs, nil, nil, nil, nil) }, name+"5 with obs3 == nil")
	shouldPanic(t, func() { _ = f5(obs, obs, obs, nil, nil, nil) }, name+"5 with obs4 == nil")
	shouldPanic(t, func() { _ = f5(obs, obs, obs, obs, nil, nil) }, name+"5 with obs5 == nil")
	shouldPanic(t, func() { _ = f5(obs, obs, obs, obs, obs, nil) }, name+"5 with proj == nil")

	shouldPanic(t, func() { _ = f6(nil, nil, nil, nil, nil, nil, nil) }, name+"6 with obs1 == nil")
	shouldPanic(t, func() { _ = f6(obs, nil, nil, nil, nil, nil, nil) }, name+"6 with obs2 == nil")
	shouldPanic(t, func() { _ = f6(obs, obs, nil, nil, nil, nil, nil) }, name+"6 with obs3 == nil")
	shouldPanic(t, func() { _ = f6(obs, obs, obs, nil, nil, nil, nil) }, name+"6 with obs4 == nil")
	shouldPanic(t, func() { _ = f6(obs, obs, obs, obs, nil, nil, nil) }, name+"6 with obs5 == nil")
	shouldPanic(t, func() { _ = f6(obs, obs, obs, obs, obs, nil, nil) }, name+"6 with obs6 == nil")
	shouldPanic(t, func() { _ = f6(obs, obs, obs, obs, obs, obs, nil) }, name+"6 with proj == nil")

	shouldPanic(t, func() { _ = f7(nil, nil, nil, nil, nil, nil, nil, nil) }, name+"7 with obs1 == nil")
	shouldPanic(t, func() { _ = f7(obs, nil, nil, nil, nil, nil, nil, nil) }, name+"7 with obs2 == nil")
	shouldPanic(t, func() { _ = f7(obs, obs, nil, nil, nil, nil, nil, nil) }, name+"7 with obs3 == nil")
	shouldPanic(t, func() { _ = f7(obs, obs, obs, nil, nil, nil, nil, nil) }, name+"7 with obs4 == nil")
	shouldPanic(t, func() { _ = f7(obs, obs, obs, obs, nil, nil, nil, nil) }, name+"7 with obs5 == nil")
	shouldPanic(t, func() { _ = f7(obs, obs, obs, obs, obs, nil, nil, nil) }, name+"7 with obs6 == nil")
	shouldPanic(t, func() { _ = f7(obs, obs, obs, obs, obs, obs, nil, nil) }, name+"7 with obs7 == nil")
	shouldPanic(t, func() { _ = f7(obs, obs, obs, obs, obs, obs, obs, nil) }, name+"7 with proj == nil")

	shouldPanic(t, func() { _ = f8(nil, nil, nil, nil, nil, nil, nil, nil, nil) }, name+"8 with obs1 == nil")
	shouldPanic(t, func() { _ = f8(obs, nil, nil, nil, nil, nil, nil, nil, nil) }, name+"8 with obs2 == nil")
	shouldPanic(t, func() { _ = f8(obs, obs, nil, nil, nil, nil, nil, nil, nil) }, name+"8 with obs3 == nil")
	shouldPanic(t, func() { _ = f8(obs, obs, obs, nil, nil, nil, nil, nil, nil) }, name+"8 with obs4 == nil")
	shouldPanic(t, func() { _ = f8(obs, obs, obs, obs, nil, nil, nil, nil, nil) }, name+"8 with obs5 == nil")
	shouldPanic(t, func() { _ = f8(obs, obs, obs, obs, obs, nil, nil, nil, nil) }, name+"8 with obs6 == nil")
	shouldPanic(t, func() { _ = f8(obs, obs, obs, obs, obs, obs, nil, nil, nil) }, name+"8 with obs7 == nil")
	shouldPanic(t, func() { _ = f8(obs, obs, obs, obs, obs, obs, obs, nil, nil) }, name+"8 with obs8 == nil")
	shouldPanic(t, func() { _ = f8(obs, obs, obs, obs, obs, obs, obs, obs, nil) }, name+"8 with proj == nil")

	shouldPanic(t, func() { _ = f9(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil) }, name+"9 with obs1 == nil")
	shouldPanic(t, func() { _ = f9(obs, nil, nil, nil, nil, nil, nil, nil, nil, nil) }, name+"9 with obs2 == nil")
	shouldPanic(t, func() { _ = f9(obs, obs, nil, nil, nil, nil, nil, nil, nil, nil) }, name+"9 with obs3 == nil")
	shouldPanic(t, func() { _ = f9(obs, obs, obs, nil, nil, nil, nil, nil, nil, nil) }, name+"9 with obs4 == nil")
	shouldPanic(t, func() { _ = f9(obs, obs, obs, obs, nil, nil, nil, nil, nil, nil) }, name+"9 with obs5 == nil")
	shouldPanic(t, func() { _ = f9(obs, obs, obs, obs, obs, nil, nil, nil, nil, nil) }, name+"9 with obs6 == nil")
	shouldPanic(t, func() { _ = f9(obs, obs, obs, obs, obs, obs, nil, nil, nil, nil) }, name+"9 with obs7 == nil")
	shouldPanic(t, func() { _ = f9(obs, obs, obs, obs, obs, obs, obs, nil, nil, nil) }, name+"9 with obs8 == nil")
	shouldPanic(t, func() { _ = f9(obs, obs, obs, obs, obs, obs, obs, obs, nil, nil) }, name+"9 with obs9 == nil")
	shouldPanic(t, func() { _ = f9(obs, obs, obs, obs, obs, obs, obs, obs, obs, nil) }, name+"9 with proj == nil")
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
