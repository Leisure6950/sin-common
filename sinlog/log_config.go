package sinlog

import "go.uber.org/zap/zapcore"

type logLevel string

var logLevels = struct {
	Debug logLevel
	Info  logLevel
	Warn  logLevel
	Error logLevel
}{
	Debug: "debug",
	Info:  "info",
	Warn:  "warn",
	Error: "error",
}

func (l logLevel) GetZapValue() zapcore.Level {
	switch l {
	case logLevels.Debug:
		return zapcore.DebugLevel
	case logLevels.Info:
		return zapcore.InfoLevel
	case logLevels.Warn:
		return zapcore.WarnLevel
	case logLevels.Error:
		return zapcore.ErrorLevel
	}
	return zapcore.InfoLevel
}

type logRotate string

var logRotates = struct {
	Hour  logRotate
	Day   logRotate
	Month logRotate
}{
	Hour:  "hour",
	Day:   "day",
	Month: "month",
}

func (r logRotate) GetRollingFormat() rollingFormat {
	switch r {
	case logRotates.Hour:
		return rollingFormats.Hourly
	case logRotates.Day:
		return rollingFormats.Daily
	case logRotates.Month:
		return rollingFormats.Monthly
	}
	return rollingFormats.Daily
}

//
// 日志配置
//
type Config struct {
	Path   string    `toml:"path" json:"path"`
	Prefix string    `toml:"prefix" json:"prefix"`
	Level  logLevel  `toml:"level" json:"level"`
	Rotate logRotate `toml:"rotate" json:"rotate"`
}

var _globalConfig = Config{
	Path:   "./logs",
	Prefix: "",
	Level:  logLevels.Info,
	Rotate: logRotates.Day,
}
