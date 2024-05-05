package utils

import (
	"go.uber.org/zap"
	"log"
)

var (
	RootLogger *zap.Logger
	Logger     *zap.SugaredLogger
)

func init() {
	var err error
	RootLogger, err = zap.NewProduction()
	//fixme 这里的错误需要处理吗？
	defer RootLogger.Sync()
	if err != nil {
		log.Fatalln("Init logger fail :", err)
	}
	Logger = RootLogger.Sugar()
}
