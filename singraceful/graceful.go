package singraceful

import "sync"

type gracefuleDispatcher struct {
	wg sync.WaitGroup
	//lock     sync.RWMutex
	disabled bool
}

var gracefulDisp gracefuleDispatcher

func destroy() {
	//gracefulDisp.lock.Lock()
	gracefulDisp.disabled = true
	//gracefulDisp.lock.Unlock()
	gracefulDisp.wg.Wait()
}

// return：是否已
func register(num int) {
	//gracefulDisp.lock.RLock()
	//defer gracefulDisp.lock.RUnlock()
	//if gracefulDisp.disabled {
	//	return false
	//}
	gracefulDisp.wg.Add(num)
	//return gracefulDisp.disabled
}

//func registerForce(num int) bool {
//	gracefulDisp.lock.Lock()
//	defer gracefulDisp.lock.Unlock()
//	gracefulDisp.wg.Add(num)
//	return true
//}

func done() {
	gracefulDisp.wg.Done()
}
func isDisabled() bool {
	return gracefulDisp.disabled
}
