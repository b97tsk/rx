package rx

// Serialize returns an Observer that passes incoming emissions to o
// in a mutually exclusive way.
// Serialize also returns a copy of c that will be cancelled when o is
// about to receive a notification of error or completion.
func Serialize[T any](c Context, o Observer[T]) (Context, Observer[T]) {
	c, cancel := c.WithCancel()
	u := new(unicast[T])
	u.Subscribe(c, o.DoOnTermination(cancel))
	return c, u.Emit
}
