package logger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var ZapLog *zap.SugaredLogger

func InitLogger() {
	var err error
	var logger *zap.Logger

	today := time.Now().Format("20060102")
	logFile := filepath.Join(viper.GetString("app.homePath"), fmt.Sprintf("gomonitor-%s.log", today))
	if _, err := os.Stat(logFile); errors.Is(err, os.ErrNotExist) {
		os.OpenFile(logFile, os.O_RDONLY|os.O_CREATE, 0777)
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
