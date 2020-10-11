package sinlog

import (
	"go.uber.org/zap/zapcore"
	"time"
)

func milliSecondTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}
var defaultEncoderConfig = zapcore.EncoderConfig{
	CallerKey:      "caller",
	StacktraceKey:  "stack",
	LineEnding:     zapcore.DefaultLineEnding,
	TimeKey:        "time",
	MessageKey:     "msg",
	LevelKey:       "level",
	NameKey:        "logger",
	EncodeCaller:   zapcore.ShortCallerEncoder,
	EncodeLevel:    zapcore.CapitalColorLevelEncoder,
	EncodeTime:     milliSecondTimeEncoder,
	EncodeDuration: zapcore.StringDurationEncoder,
	EncodeName:     zapcore.FullNameEncoder,
}
