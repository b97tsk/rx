package rx

// Option is the basic type of options.
type Option interface {
	thisIsAnOption()
}

type (
	compareOption struct {
		Value func(interface{}, interface{}) bool
	}

	concurrentOption struct {
		Value int
	}

	defaultValueOption struct {
		Value interface{}
	}

	keySelectorOption struct {
		Value func(interface{}) interface{}
	}

	predicateOption struct {
		Value func(interface{}, int) bool
	}

	projectOption struct {
		Value func(interface{}, int) interface{}
	}

	resultSelectorOption struct {
		Value func(interface{}, interface{}, int, int) interface{}
	}

	seedOption struct {
		Value interface{}
	}

	schedulerOption struct {
		Value Scheduler
	}
)

func (compareOption) thisIsAnOption()        {}
func (concurrentOption) thisIsAnOption()     {}
func (defaultValueOption) thisIsAnOption()   {}
func (keySelectorOption) thisIsAnOption()    {}
func (predicateOption) thisIsAnOption()      {}
func (projectOption) thisIsAnOption()        {}
func (resultSelectorOption) thisIsAnOption() {}
func (seedOption) thisIsAnOption()           {}
func (schedulerOption) thisIsAnOption()      {}

// An OptionSlice is a slice of options. It has many useful methods that lets
// you easily create an option.
type OptionSlice []Option

// WithCompare creates a compare option, appends it to the end of the slice,
// and returns the updated slice.
func (s OptionSlice) WithCompare(v func(interface{}, interface{}) bool) OptionSlice {
	return append(s, compareOption{v})
}

// WithConcurrent creates a concurrent option, appends it to the end of the slice,
// and returns the updated slice.
func (s OptionSlice) WithConcurrent(v int) OptionSlice {
	return append(s, concurrentOption{v})
}

// WithDefaultValue creates a default-value option, appends it to the end of the slice,
// and returns the updated slice.
func (s OptionSlice) WithDefaultValue(v interface{}) OptionSlice {
	return append(s, defaultValueOption{v})
}

// WithKeySelector creates a key-selector option, appends it to the end of the slice,
// and returns the updated slice.
func (s OptionSlice) WithKeySelector(v func(interface{}) interface{}) OptionSlice {
	return append(s, keySelectorOption{v})
}

// WithPredicate creates a predicate option, appends it to the end of the slice,
// and returns the updated slice.
func (s OptionSlice) WithPredicate(v func(interface{}, int) bool) OptionSlice {
	return append(s, predicateOption{v})
}

// WithProject creates a project option, appends it to the end of the slice,
// and returns the updated slice.
func (s OptionSlice) WithProject(v func(interface{}, int) interface{}) OptionSlice {
	return append(s, projectOption{v})
}

// WithResultSelector creates a result-selector option, appends it to the end
// of the slice, and returns the updated slice.
func (s OptionSlice) WithResultSelector(v func(interface{}, interface{}, int, int) interface{}) OptionSlice {
	return append(s, resultSelectorOption{v})
}

// WithSeed creates a seed option, appends it to the end of the slice, and
// returns the updated slice.
func (s OptionSlice) WithSeed(v interface{}) OptionSlice {
	return append(s, seedOption{v})
}

// WithScheduler creates a scheduler option, appends it to the end of the
// slice, and returns the updated slice.
func (s OptionSlice) WithScheduler(v Scheduler) OptionSlice {
	return append(s, schedulerOption{v})
}

// WithCompare returns a new OptionSlice that contains a new compare option.
func WithCompare(v func(interface{}, interface{}) bool) OptionSlice {
	return OptionSlice(nil).WithCompare(v)
}

// WithConcurrent returns a new OptionSlice that contains a new concurrent option.
func WithConcurrent(v int) OptionSlice {
	return OptionSlice(nil).WithConcurrent(v)
}

// WithDefaultValue returns a new OptionSlice that contains a new default-value option.
func WithDefaultValue(v interface{}) OptionSlice {
	return OptionSlice(nil).WithDefaultValue(v)
}

// WithKeySelector returns a new OptionSlice that contains a new key-selector option.
func WithKeySelector(v func(interface{}) interface{}) OptionSlice {
	return OptionSlice(nil).WithKeySelector(v)
}

// WithPredicate returns a new OptionSlice that contains a new predicate option.
func WithPredicate(v func(interface{}, int) bool) OptionSlice {
	return OptionSlice(nil).WithPredicate(v)
}

// WithProject returns a new OptionSlice that contains a new project option.
func WithProject(v func(interface{}, int) interface{}) OptionSlice {
	return OptionSlice(nil).WithProject(v)
}

// WithResultSelector returns a new OptionSlice that contains a new result-selector option.
func WithResultSelector(v func(interface{}, interface{}, int, int) interface{}) OptionSlice {
	return OptionSlice(nil).WithResultSelector(v)
}

// WithSeed returns a new OptionSlice that contains a new seed option.
func WithSeed(v interface{}) OptionSlice {
	return OptionSlice(nil).WithSeed(v)
}

// WithScheduler returns a new OptionSlice that contains a new scheduler option.
func WithScheduler(v Scheduler) OptionSlice {
	return OptionSlice(nil).WithScheduler(v)
}
