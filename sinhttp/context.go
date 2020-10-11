package sinhttp

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sin-z/sin-common/sinlog"
	"go.uber.org/zap"
	"sync"
	"time"
)

type env struct {
	ctxMap map[context.Context]*session
	lock   sync.RWMutex
}

var _env *env

func init() {
	_env = &env{
		ctxMap: make(map[context.Context]*session),
	}
}

func Session(ctx context.Context) *session {
	return _env.Session(ctx)
}
func HasSession(ctx context.Context) bool {
	return _env.Exist(ctx)
}

func (e *env) Session(ctx context.Context) *session {
	var sess *session
	e.lock.Lock()
	sess, ok := e.ctxMap[ctx]
	if !ok {
		sessDone := make(chan struct{}, 0)
		// generate session
		sess = &session{
			ctx:   ctx,
			env:   e,
			cache: make(map[string]interface{}),
			done:  sessDone,
		}
		// save to map
		e.ctxMap[ctx] = sess
		// 默认不过期
		//sess.WithTimeout(time.Minute * 3)
		// listen done
		go func() {
			select {
			case <-ctx.Done():
				sinlog.For(ctx, zap.String("func", "Session")).Infow("request handle done")
				sess.cancel()
				return
			case <-sessDone:
				sinlog.For(ctx, zap.String("func", "Session")).Infow("request handle timeout")
				sess.cancel()
				return
			}
		}()
	}
	e.lock.Unlock()
	return sess
}
func (e *env) Exist(ctx context.Context) bool {
	e.lock.RLock()
	_, ok := e.ctxMap[ctx]
	e.lock.RUnlock()
	return ok
}

type session struct {
	ctx   context.Context
	env   *env
	done  chan struct{}
	timer *time.Timer

	cache map[string]interface{}
	lock  sync.RWMutex
}

func (s *session) Add(k string, v interface{}) error {
	s.lock.Lock()
	if oldV, ok := s.cache[k]; ok {
		if oldV == v {
			s.lock.Unlock()
			return nil
		}
		s.lock.Unlock()
		return errors.Errorf("key repeated")
	}
	s.cache[k] = v
	s.lock.Unlock()
	return nil
}
func (s *session) Set(k string, v interface{}) {
	s.lock.Lock()
	s.cache[k] = v
	s.lock.Unlock()
}

func (s *session) Get(k string) (interface{}, bool) {
	s.lock.RLock()
	v, ok := s.cache[k]
	s.lock.RUnlock()
	if !ok {
		return nil, false
	}
	return v, ok
}

func (s *session) WithTimeout(d time.Duration) *session {
	s.lock.Lock()
	if s.timer != nil {
		s.timer.Stop()
		s.timer = nil
	}
	// 过期时间>0，设置过期处理。<=0则永久不过期
	if d > 0 {
		s.timer = time.AfterFunc(d, func() {
			s.cancel()
		})
	}
	s.lock.Unlock()
	return s
}
func (s *session) cancel() {
	s.lock.Lock()
	if s.timer != nil {
		s.timer.Stop()
		s.timer = nil
	}
	if s.done != nil {
		tmpDone := s.done
		s.done = nil
		close(tmpDone)
	}
	delete(s.env.ctxMap, s.ctx)
	s.lock.Unlock()
}
