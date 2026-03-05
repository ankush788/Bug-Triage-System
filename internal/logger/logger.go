package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func Init() {
    var err error
    Log, err = zap.NewProduction()
    if err != nil {
        panic(err)
    }
    zap.ReplaceGlobals(Log)
}

func Sync() {
    _ = Log.Sync()
}