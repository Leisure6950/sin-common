package sinlog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
)

//
// 日志对象
//
type logger struct {
	*zap.SugaredLogger

	path    string
	rolling rollingFormat

	rollingFiles []io.Writer

	loglevel zap.AtomicLevel
	prefix   string

	encoderCfg zapcore.EncoderConfig
}

func newLogger() *logger {
	cfg := defaultEncoderConfig
	lvl := zap.NewAtomicLevelAt(_globalConfig.Level.GetZapValue())

	log := &logger{
		SugaredLogger: zap.New(
			zapcore.NewCore(
				newConsoleEncoder(&cfg, false),
				//zapcore.NewConsoleEncoder(cfg),
				//zapcore.NewJSONEncoder(cfg),
				zapcore.Lock(os.Stderr), lvl),
		).WithOptions(
			zap.AddCaller(),
			zap.AddCallerSkip(1),
		).Sugar(),
		path:         "",
		rolling:      "",
		rollingFiles: nil,
		loglevel:     lvl,
		prefix:       _globalConfig.Prefix,
		encoderCfg:   cfg,
	}

	return log
}

func (l *logger) closeFiles() {
	for _, w := range l.rollingFiles {
		r, ok := w.(*rollingFile)
		if ok {
			r.Close()
		}
	}
	l.rollingFiles = nil
}
func (l *logger) refreshRotate() {
	for _, w := range l.rollingFiles {
		r, ok := w.(*rollingFile)
		if ok {
			r.SetRolling(l.rolling)
		}
	}
}

// Up to tow log files
// full_log、error_log
func (l *logger) setOutputPath(path string) error {
	if l.path == path {
		return nil
	}
	l.closeFiles()
	l.path = path
	fullFile, err := newRollingFile(path+"/full.log", l.rolling)
	if err != nil {
		return err
	}
	//infoFile, err := newRollingFile(path+"/info.log", l.rolling)
	//if err != nil {
	//	return err
	//}
	errorFile, err := newRollingFile(path+"/error.log", l.rolling)
	if err != nil {
		return err
	}
	fullLogEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		//if l.loglevel.Level() >= zapcore.InfoLevel {
		//	return false
		//}
		//return l.loglevel.Enabled(lvl)
		return true
	})
	errorLogEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})
	cfg := l.encoderCfg
	//cfg.LevelKey = "lvl"
	//cfg.MessageKey = "topic"
	core := zapcore.NewTee(
		zapcore.NewCore(newConsoleEncoder(&cfg, false), fullFile, fullLogEnabler),
		zapcore.NewCore(newConsoleEncoder(&cfg, false), errorFile, errorLogEnabler),
	)
	//ore := zapcore.NewTee(
	//	zapcore.NewCore(zapcore.NewConsoleEncoder(cfg), fullFile, fullLogEnabler),
	//	zapcore.NewCore(newConsoleEncoder(&cfg, false), errorFile, errorLogEnabler),
	//)

	//core := zapcore.NewCore(
	//	newConsoleEncoder(&cfg, false),
	//	zapcore.Lock(os.Stderr), l.loglevel)

	l.rollingFiles = []io.Writer{fullFile, errorFile}

	l.SugaredLogger = zap.New(core).WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
	l.SugaredLogger.Named(l.prefix)
	return nil
}
func (l *logger) setLevel(level logLevel) {
	l.loglevel.SetLevel(level.GetZapValue())
}
func (l *logger) setRolling(rotate logRotate) {
	l.rolling = rotate.GetRollingFormat()
	l.refreshRotate()
}
func (l *logger) setLogPrefix(prefix string) {
	l.prefix = prefix
	l.SugaredLogger.Named(prefix)
}
func (l *logger) With(args ...interface{}) *logger {
	return &logger{SugaredLogger: l.SugaredLogger.With(args...).Desugar().WithOptions(zap.AddCallerSkip(0)).Sugar()}
}
