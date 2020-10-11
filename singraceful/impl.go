package singraceful

import (
	"github.com/pkg/errors"
	"github.com/sin-z/sin-common/sinlog"
	"github.com/sin-z/sin-common/sinprocess"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

type serverRunner struct {
	once     sync.Once
	lock     sync.Mutex
	num      int
	lastTime time.Time
	srvDone  chan struct{}
}

var r = &serverRunner{
	srvDone: make(chan struct{}),
}

func (r *serverRunner) startServerAndWait(launcher func() error) (err error) {
	err = func() error {
		r.lock.Lock()
		defer r.lock.Unlock()
		if r.num > 0 {
			return errors.Errorf("Start can only be called one time")
		}
		r.num++
		if r.srvDone == nil {
			return errors.Errorf("Service already stopped")
		}
		return nil
	}()
	if err != nil {
		return
	}
	// 主处理逻辑退出
	if launcher != nil {
		go func() {
			err := launcher()
			if err != nil && err != http.ErrServerClosed {
				sinlog.With(zap.String("func", "startAndWait")).Errorw("http server err", zap.Error(err))
				sinlog.Sync()
			}
			r.lock.Lock()
			if r.srvDone != nil {
				close(r.srvDone)
				r.srvDone = nil
			}
			r.lock.Unlock()
		}()
	}

	// 系统停服信号
	sinprocess.ExitCallback.Register(func(err error) {
		r.lock.Lock()
		if r.srvDone != nil {
			close(r.srvDone)
			r.srvDone = nil
		}
		r.lock.Unlock()
	})

	// 等待系统停止信号
	r.lock.Lock()
	doneAct := r.srvDone
	r.lock.Unlock()
	if doneAct != nil {
		select {
		case <-doneAct:
		}
	} else {
		sinlog.With(zap.String("func", "startServerAndWait")).Warnw("The service maybe not init or already stopped")
	}

	// 待所有消息处理完成之后，退出系统
	register(1)
	go func() {
		defer done()
		for { // 2秒内没有新任务计数，则继续退出
			time.Sleep(time.Millisecond * 100)
			r.lock.Lock()
			if time.Now().UnixNano()-r.lastTime.UnixNano() > 1000*1000*1000*2 {
				r.lock.Unlock()
				break
			}
			r.lock.Unlock()
		}
	}()
	// 等待任务处理完成
	destroy()
	return
}

type asyncRunner struct {
	wg  sync.WaitGroup
	err error
}

func (ar *asyncRunner) Wait() *asyncRunner {
	ar.wg.Wait()
	return ar
}
func (ar *asyncRunner) Error() error {
	ar.wg.Wait()
	return ar.err
}
func (ar *asyncRunner) Go(f func() error) *asyncRunner {
	r.lock.Lock()
	r.lastTime = time.Now()
	r.lock.Unlock()
	return ar.async(f)
}

/**

 */
func (ar *asyncRunner) async(f func() error) *asyncRunner {
	register(1)
	ar.wg.Add(1)
	go func() {
		defer func() {
			ar.wg.Done()
			done()
		}()
		ar.err = f()
	}()
	return ar
}
