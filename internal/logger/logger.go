package logger

import (
	"go.uber.org/zap"
)

// zap is a high-performance structured logging library for Go created by Uber
// it is faster because
// 1)it's internal implentation
//2) structure : fmt.Println log have no structure but  logs have (json) so, it easily parse by machine

var Log *zap.Logger

func Init() {
    var err error
    // Log, err = zap.NewProduction() //This creates a production configuration logger.
    // 1) ignore debug level log 
    // 2) write all log in json fomat (machine friendly format)
    // logger.Info("user created", zap.String("user_id", "123"))
    // converted into 
    //{
    //   "level": "info",
    //   "msg": "user created",
    //   "user_id": "123"
    // }

    Log, err = zap.NewDevelopment() //This creates a development configuration logger.
    //1)  it include all level like (debug , Info , Warn, Error, Fatal)
    //2) it write data in human friendly format (Still structed)  
    //ex: 2026-03-05T10:20:00.123+0530 DEBUG main.go:10 {debug_message}
    if err != nil {
        panic(err)
    }

}

// zap write log in buffer and batch processing to disk 
// sync() -> Flush remaining buffered logs to disk/output.
// without it some logs are not properly written in disk
func Sync() {
    _ = Log.Sync()
}
