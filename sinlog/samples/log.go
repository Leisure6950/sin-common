package main

import (
	"context"
	"github.com/sin-z/sin-common/sinlog"
	"go.uber.org/zap"
)

func main() {
	sinlog.Init(sinlog.Config{
		Path:   "./logs",
		Prefix: "test",
		Level:  "debug",
		Rotate: "hour",
	})

	log := sinlog.For(context.Background(), zap.String("key1", "1"))

	log.Debugw("I am debug", "k", "v")
	log.Infow("I am info", "k", "v")
	log.Warnw("I am warn", "k", "v")
	log.Errorw("I am error", "k", "v")

	log = log.With(zap.String("key2", "2"))

	log.Debugw("I am debug", "k", "v")
	log.Infow("I am info", "k", "v")
	log.Warnw("I am warn", "k", "v")
	log.Errorw("I am error", "k", "v")

	log.Sync()
}
