package sinutil

import (
	"sync"
	"sync/atomic"
)

type OnceSucceed struct {
	m       sync.Mutex
	succeed uint32
}

func (o *OnceSucceed) Do(f func() error) (err error) {
	if atomic.LoadUint32(&o.succeed) == 1 {
		return nil
	}
	// Slow-path.
	o.m.Lock()
	defer o.m.Unlock()
	if o.succeed == 0 {
		if err = f(); err == nil {
			atomic.StoreUint32(&o.succeed, 1)
		}
	}
	return err
}
