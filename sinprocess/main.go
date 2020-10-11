package sinprocess

import (
	"context"
	"github.com/sin-z/sin-common/sinutil"
	"os/signal"
	"sync"
	"syscall"
)

type processDoneFunc func(err error)

//type processExitAction int
//
//var ProcessExitActions = struct {
//	Exit      processExitAction // 继续推出
//	Interrupt processExitAction // 中断
//}{}

var processer = struct {
	ctx       context.Context
	observers []processDoneFunc
	lock      sync.RWMutex
	once      sync.Once
}{
	ctx: sinutil.WithSignals(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL),
}

func listen() {
	go func() {
		Wait()
		processer.lock.Lock()
		observers := processer.observers
		processer.lock.Unlock()
		for _, obs := range observers {
			obs(nil)
		}
	}()
}

type exitCallback struct{}

var ExitCallback = &exitCallback{}

func (*exitCallback) Register(doneFunc processDoneFunc) {
	processer.once.Do(func() {
		listen()
	})
	processer.lock.Lock()
	defer processer.lock.Unlock()
	processer.observers = append(processer.observers, doneFunc)
}
func Wait() {
	processer.lock.Lock()
	ctx := processer.ctx
	processer.lock.Unlock()
	if ctx == nil {
		return
	}
	select {
	case <-ctx.Done():
	}
}
func Close() {
	// 发送可监控的系统推出信号
	signal.Ignore(syscall.SIGTERM)
	Wait()
}
