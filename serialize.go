package rx

// Serialize returns an Observer that passes incoming emissions to sink
// in a mutually exclusive way.
// Serialize also returns a copy of c that will be cancelled when sink is
// about to receive a notification of error or completion.
func Serialize[T any](c Context, sink Observer[T]) (Context, Observer[T]) {
	c, cancel := c.WithCancel()
	u := new(unicast[T])
	u.subscribe(c, sink.OnLastNotification(cancel))
	return c, u.emit
}
