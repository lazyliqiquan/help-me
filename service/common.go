package service

import (
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func Init(loggerInstance *zap.SugaredLogger) {
	Logger = loggerInstance
}
