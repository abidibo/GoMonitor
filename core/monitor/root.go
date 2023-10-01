package monitor

import (
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	"github.com/abidibo/gomonitor/core/utils"
	"github.com/abidibo/gomonitor/logger"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/spf13/viper"
)

type ProcessLog struct {
	Name          string
	CpuPercent    float64
	MemoryPercent float32
}

func RunAsRoot() {
	// now
	startTime := time.Now()
	// every logIntervalMinutes minutes we log
	logIntervalMinutes := viper.GetInt("app.logIntervalMinutes")

	// remove old data and log files
	cleanup()

	// keep program running
	for {

		// get current user
		currentUser, err := utils.GetCurrentUser()
		if err != nil {
			logger.ZapLog.Error("Cannot get current user")
		} else {
			logger.ZapLog.Debug("Current user ", currentUser)
		}

		timeScreenLimit, err := utils.GetScreenTimeLimitMinutes(currentUser)

		logger.ZapLog.Debug("Retrieving current user processes")
		logProcesses := make([]ProcessLog, 0)
		processes, _ := process.Processes()
		for _, process := range processes {
			name, _ := process.Name()
			cpuUsage, _ := process.CPUPercent()
			memoryUsage, _ := process.MemoryPercent()
			u, _ := process.Username()
			if u == currentUser && (cpuUsage > 0 || memoryUsage > 0) {
				// add process to list if it belongs to current user and is not in the list
				log := ProcessLog{Name: name, CpuPercent: cpuUsage, MemoryPercent: memoryUsage}
				logProcesses = append(logProcesses, log)
			}
		}
		sort.Slice(logProcesses, func(i, j int) bool {
			return logProcesses[i].CpuPercent > logProcesses[j].CpuPercent
		})

		// insert log into DB
		elapsed := int(math.Round(time.Since(startTime).Minutes()))
		startTime = time.Now() // reset start time
		logId, err := utils.InsertLog(currentUser, elapsed)
		if err == nil {
			for _, p := range logProcesses {
				utils.InsertProcessLog(logId, p.Name, p.CpuPercent, p.MemoryPercent)
			}
		}

		// get sum of partial_time_minutes for current day and user
		totalMinutes, err := utils.GetTotalTodayTimeMinutes(currentUser)
		if err == nil {
			// check if total time usage has exceeded the limit
			if timeScreenLimit > 0 {
				if totalMinutes > timeScreenLimit {
					logger.ZapLog.Info("Limit exceeded for user ", currentUser, " ", totalMinutes, " minutes")
					// logout user
					err = utils.LogoutUser(currentUser)
					if err != nil {
						// try to shutdown pc
						utils.Shutdown()
					}
				}
			}
		}

		time.Sleep(time.Duration(logIntervalMinutes) * time.Minute)
	}
}

func cleanup() {
	retainPeriodDays := viper.GetInt("app.retainPeriodDays")
	now := time.Now()
	limit := now.AddDate(0, 0, -1*retainPeriodDays)

	affected, err := utils.DeleteProcessData(limit.Format("2006-01-02"))
	// clean db process data
	if err == nil {
		logger.ZapLog.Info("Removed process data before ", limit, " affected ", affected)
	}

	affected, err = utils.DeleteLogData(limit.Format("2006-01-02"))
	// clean db log data
	if err == nil {
		logger.ZapLog.Info("Removed log data before ", limit, " affected ", affected)
	}

	homePath := viper.GetString("app.homePath")
	filepath.Walk(homePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.ZapLog.Error("Error walking path ", path)
		}
		re := regexp.MustCompile(`gomonitor-(\d+).log`)
		match := re.FindStringSubmatch(info.Name())
		if len(match) > 0 {
			dateLayout := "20060102"
			d, err := time.Parse(dateLayout, match[1])
			if err == nil {
				if d.Before(limit) {
					err := os.Remove(path)
					if err != nil {
						logger.ZapLog.Error("Error removing file ", path)
					}
				}
			}
		}
		return nil
	})
}
