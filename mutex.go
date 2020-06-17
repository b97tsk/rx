package rx

// Mutex creates an Observer that passes all emissions to a specified Observer
// in a mutually exclusive way.
func Mutex(sink Observer) Observer {
	cx := make(chan Observer, 1)
	cx <- sink
	return func(t Notification) {
		if sink, ok := <-cx; ok {
			switch {
			case t.HasValue:
				sink(t)
				cx <- sink
			default:
				close(cx)
				sink(t)
			}
		}
	}
}
