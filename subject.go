package rx

// A Subject is both an Observable and an Observer.
// Observers that subscribed to Subject's Observable part may receive
// emissions from Subject's Observer part.
type Subject[T any] struct {
	Observable[T]
	Observer[T]
}
