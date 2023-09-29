package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/abidibo/gomonitor/core/monitor"
	"github.com/abidibo/gomonitor/core/stats"
	"github.com/abidibo/gomonitor/core/utils"
	"github.com/abidibo/gomonitor/db"
	"github.com/abidibo/gomonitor/logger"
	"github.com/akamensky/argparse"
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
	// Check arguments
	parser := argparse.NewParser("GoMonitor", "Simple parental control application")

	// Add top level command `monitor`
	monitorCmd := parser.NewCommand("monitor", "Launch monitor")

	// Add top level command `stats`
	statsCmd := parser.NewCommand("stats", "View statistics")

	// if not root user can only view his stats
	statsUserDefault := ""
	statsUser := &statsUserDefault
	if utils.IsRunningAsRoot() {
		statsUser = statsCmd.String("u", "user", &argparse.Options{
			Required: true,
			Help:     "the user to view statistics",
		})
	}
	statsDate := statsCmd.String("d", "date", &argparse.Options{
		Required: false,
		Help:     "the date to view statistics",
		Default:  time.Now().Format("2006-01-02"),
		Validate: func(args []string) error {
			dateLayout := "2006-01-02"
			_, err := time.Parse(dateLayout, args[0])
			if err != nil {
				return errors.New("enter a valid date: YYYY-MM-DD")
			}
			return nil
		},
	})

	// Parse command line arguments and in case of any error print error and help information
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	// Init database
	db.InitDatabase()

	if monitorCmd.Happened() {
		// statistics
		monitor.Run()
	} else if statsCmd.Happened() {
		// monitor
		stats.Run(*statsUser, *statsDate)
	}
}
