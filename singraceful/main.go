package singraceful

import (
	"github.com/sin-z/sin-common/sinlog"
	"net/http"
	"time"
)

// 运行，阻塞等待服务停止或异常退出
// 示例：
//     gaiagraceful.Run(func() error {
//         return RunMsgHandle()
//     })
func Run(launcher func() error) error {
	// 启动运行监控
	err := r.startServerAndWait(launcher)
	if err != nil {
		return err
	}
	go func() {
		addr := "0.0.0.0:6060"
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			sinlog.Log().Warnf("start pprof server (%s) error, err:%v", addr, err)
		}
	}()
	return nil
}

// 执行异步任务，更为复杂的用法，可以使用gaiaasync，支持graceful
func Go(f func() error) *asyncRunner {
	return new(asyncRunner).Go(f)
}

// 执行同步任务，更为复杂的用法，可以使用gaiaasync，支持graceful
func Sync(f func() error) error {
	register(1)
	defer done()
	r.lock.Lock()
	r.lastTime = time.Now()
	r.lock.Unlock()
	return f()
}

// 循环定时调用
func Ticker(duration time.Duration, f interface{}, args ...interface{}) {
	go func() {
		defer done()
		register(1)
		r.lock.Lock()
		srvDone := r.srvDone
		r.lock.Unlock()
		if srvDone == nil {
			return
		}
		ticker := time.NewTicker(duration)
		for {
			select {
			case <-ticker.C:
				if isDisabled() {
					return
				}
				if len(args) > 1 {
					f.(func(...interface{}))(args...)
				} else if len(args) == 1 {
					f.(func(interface{}))(args[0])
				} else {
					f.(func())()
				}
			case <-r.srvDone:
				ticker.Stop()
				return
			}
		}
	}()
}

// 定时调用，仅调用一次。循环调用请使用Ticker
func Timer(duration time.Duration, f interface{}, args ...interface{}) {
	go func() {
		defer done()
		register(1)
		timer := time.NewTimer(duration)
		select {
		case <-timer.C:
			if isDisabled() {
				return
			}
			if len(args) > 1 {
				f.(func(...interface{}))(args...)
			} else if len(args) == 1 {
				f.(func(interface{}))(args[0])
			} else {
				f.(func())()
			}
		case <-r.srvDone:
			timer.Stop()
			return
		}
	}()
}

// 循环调用
func Loop(f interface{}, args ...interface{}) {
	go func() {
		defer done()
		register(1)
		for {
			if isDisabled() {
				return
			}
			if len(args) > 1 {
				f.(func(...interface{}))(args...)
			} else if len(args) == 1 {
				f.(func(interface{}))(args[0])
			} else {
				f.(func())()
			}
		}
	}()
}

// 添加运行任务计数（无论返回值true或false，均会将计数加1）
// return：true-正常运行；false-服务已申请停止，须终止任务
func Add(num int) {
	register(num)
	r.lock.Lock()
	r.lastTime = time.Now()
	r.lock.Unlock()
}

// 任务执行完成
func Done() {
	done()
}

// 校验服务是否已申请停止
func IsDisabled() bool {
	return isDisabled()
}
