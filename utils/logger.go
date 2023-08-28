package utils

import (
	"go.uber.org/zap"
	"san_francisco/config"
)

var Logger *zap.Logger
var SugarLogger *zap.SugaredLogger

func InitializeLogger() {
	Logger = zap.Must(zap.NewDevelopment())
	if config.Env == "PROD" {
		Logger = zap.Must(zap.NewProduction())
	}
	SugarLogger = Logger.Sugar()
}
