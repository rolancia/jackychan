package jackychan

import "sync"

type RouteBufferName string

var storage RouteBufferStorage

func Init(debug bool) {
	if !debug {
		storage = RouteBufferStorage{
			buffers: make(map[RouteBufferName]*[]RouteBuffer),
			chIn:    make(chan RouteBuffer, 2*4096),
			chOut:   make(chan routeBuffersRequest, 8),
		}
	} else {
		storage = RouteBufferStorage{
			buffers: make(map[RouteBufferName]*[]RouteBuffer),
			chIn:    make(chan RouteBuffer),
			chOut:   make(chan routeBuffersRequest),
		}
	}

	go func() {
		for {
			select {
			case b := <-storage.chIn:
				if _, ok := storage.buffers[b.Name]; !ok {
					newBuffs := make([]RouteBuffer, 0, 16)
					storage.buffers[b.Name] = &newBuffs
				}
				*storage.buffers[b.Name] = append(*storage.buffers[b.Name], b)
			case req := <-storage.chOut:
				for n, b := range storage.buffers {
					buffers := make([]RouteBuffer, len(*b))
					copy(buffers, *b)
					req.result[n] = buffers
				}
				req.wg.Done()
			}
		}
	}()
}

type RouteBufferStorage struct {
	buffers map[RouteBufferName]*[]RouteBuffer

	chIn  chan RouteBuffer
	chOut chan routeBuffersRequest
}

func Clone() map[RouteBufferName][]RouteBuffer {
	wg := sync.WaitGroup{}
	wg.Add(1)
	ret := make(map[RouteBufferName][]RouteBuffer)
	storage.chOut <- routeBuffersRequest{
		wg:     &wg,
		result: ret,
	}
	wg.Wait()
	return ret
}

func Append(buffer RouteBuffer) {
	storage.chIn <- buffer
}

type routeBuffersRequest struct {
	wg     *sync.WaitGroup
	result map[RouteBufferName][]RouteBuffer
}

func newRouteBuffer(name RouteBufferName) RouteBuffer {
	return RouteBuffer{
		Name:     name,
		Entities: make([]Entity, 8),
	}
}

type RouteBuffer struct {
	Name     RouteBufferName
	Entities []Entity
}

type Entity struct {
	Title     string
	Message   string
	Source    string
	Timestamp int64
}
