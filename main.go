package main

import (
	"fmt"
	"os"
	"time"

	"github.com/abidibo/gomonitor/core"
	"github.com/abidibo/gomonitor/db"
	"github.com/abidibo/gomonitor/logger"
	"github.com/spf13/viper"
)

func init() {
	// Read settings
	viper.SetConfigFile(fmt.Sprintf("/etc/gomonitor.json"))
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("Error reading settings file, %s", err))
	}

	// Ensure app home directory exists
	homePath := viper.GetString("app.homePath")

	err := os.MkdirAll(homePath, os.ModePerm)
	if err != nil {
		panic("Error creating the application home directory")
	}

	// Init logger
	logger.InitLogger()
	logger.ZapLog.Info("Starting gomonitor application")
}

func main() {
	// Init database
	db.InitDatabase()

	// view statistics
	if len(os.Args) > 1 && os.Args[1] == "stats" {
		if len(os.Args) < 3 {
			currentTime := time.Now()
			core.Stats(currentTime.Format("2006-01-02"))
		} else {
			core.Stats(os.Args[2])
		}
		return
	}

	// run and log
	core.Run()
}
