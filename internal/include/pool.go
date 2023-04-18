package include

import (
	"time"
)

const INC_POOL_SIZE = 1000

type incPool struct {
	pool    chan *Includes
	factory func() *Includes
}

func UseIncPool(factory func() *Includes) *incPool {
	return &incPool{
		pool:    make(chan *Includes, INC_POOL_SIZE),
		factory: factory,
	}
}

func (i *incPool) Get() *Includes {
	select {
	case inc := <-i.pool:
		return inc
	default:
		inc := i.factory()
		return inc
	}
}

func (i *incPool) Release(inc *Includes) {
	select {
	case i.pool <- inc:
	case <-time.After(time.Millisecond):
	}
}
