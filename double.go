package rx

// A Double is both an Observable and an Observer. Observers that subscribed to
// Double's Observable part may receive emissions from Double's Observer part.
type Double struct {
	Observable
	Observer
}

// A DoubleFactory is a factory function that produces Doubles.
type DoubleFactory func() Double
