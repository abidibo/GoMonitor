package logger

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var ZapLog *zap.SugaredLogger

func InitLogger() {
	var err error
	var logger *zap.Logger

	logFile := filepath.Join(viper.GetString("app.homePath"), "gomonitor.log")
	if _, err := os.Stat(logFile); errors.Is(err, os.ErrNotExist) {
		os.OpenFile(logFile, os.O_RDONLY|os.O_CREATE, 0666)
	}

	cfg := zap.NewDevelopmentConfig()

	cfg.OutputPaths = []string{
		logFile,
	}
	logger, err = cfg.Build()
	ZapLog = logger.Sugar()

	if err != nil {
		panic(err)
	}
}
