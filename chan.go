package jackychan

func Chan[T any](name RouteBufferName) Channel[T] {
	return Channel[T]{
		routeBuffer: newRouteBuffer(name),
		ch:          make(chan T),
		chBridge:    make(chan T, 256),
	}
}

func ChanAsync[T any](name RouteBufferName, size int) Channel[T] {
	return Channel[T]{
		routeBuffer: newRouteBuffer(name),
		ch:          make(chan T, size),
		chBridge:    make(chan T, 256),
	}
}

type Channel[T any] struct {
	routeBuffer RouteBuffer
	ch          chan T
	chBridge    chan T
}

func (ch Channel[T]) Push(data T) {
	ch.ch <- data
}

func (ch Channel[T]) Pop() T {
	return <-ch.ch
}

func (ch Channel[T]) Selector() chan T {
	select {
	case v := <-ch.ch:
		ch.chBridge <- v
	default:
	}
	return ch.chBridge
}
