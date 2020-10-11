package sinlog

import (
	"context"
	"github.com/sin-z/sin-common/sinprocess"
	"go.uber.org/zap"
)

const (
	_jsonDataTaskKey = "service_name"
	traceIDKey       = "trace_id"
)

// 日志
var _defaultLogger *logger

//
// 初始化日志库
//
func Init(cfg Config) {
	_globalConfig = cfg
	_defaultLogger = newLogger()
	_defaultLogger.setLevel(_globalConfig.Level)
	_defaultLogger.setLogPrefix(_globalConfig.Prefix)
	_defaultLogger.setRolling(_globalConfig.Rotate)
	_defaultLogger.setOutputPath(_globalConfig.Path)
	sinprocess.ExitCallback.Register(func(err error) {
		_defaultLogger.Sync()
	})
}

func For(ctx context.Context, args ...interface{}) *logger {
	tid := extraTraceID(ctx)
	var fields []interface{}
	if len(tid) != 0 {
		fields = make([]interface{}, 0, len(args)+2)
		fields = append(fields, traceIDKey, extraTraceID(ctx))
	} else {
		fields = make([]interface{}, 0, len(args))
	}
	fields = append(fields, args...)
	return With(fields...)
}
func With(args ...interface{}) *logger {
	return &logger{SugaredLogger: _defaultLogger.With(args...).Desugar().WithOptions(zap.AddCallerSkip(-1)).Sugar()}
}

// Sync all log data
func Sync() {
	if _defaultLogger != nil {
		_defaultLogger.Sync()
	}
}
func Log() *logger {
	return &logger{SugaredLogger: _defaultLogger.Desugar().WithOptions(zap.AddCallerSkip(-1)).Sugar()}
}

func Debug(v ...interface{}) {
	_defaultLogger.Debug(v...)
}

func Info(v ...interface{}) {
	_defaultLogger.Info(v...)
}

func Warn(v ...interface{}) {
	_defaultLogger.Warn(v...)
}

func Error(v ...interface{}) {
	_defaultLogger.Error(v...)
}

func Fatal(v ...interface{}) {
	_defaultLogger.Fatal(v...)
}

func Debugf(format string, v ...interface{}) {
	_defaultLogger.Debugf(format, v...)
}

func Infof(format string, v ...interface{}) {
	_defaultLogger.Infof(format, v...)
}

func Warnf(format string, v ...interface{}) {
	_defaultLogger.Warnf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	_defaultLogger.Errorf(format, v...)
}

func Fatalf(format string, v ...interface{}) {
	_defaultLogger.Fatalf(format, v...)
}
func Debugw(msg string, keysAndValues ...interface{}) {
	_defaultLogger.Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	_defaultLogger.Infow(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	_defaultLogger.Errorw(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	_defaultLogger.Warnw(msg, keysAndValues...)
}
