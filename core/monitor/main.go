package monitor

import (
	"fmt"
	"math"
	"os"
	"sort"
	"time"

	"github.com/abidibo/gomonitor/core/stats"
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

func Run() {
	if len(os.Args) > 1 {
		// show stats
		stats.Stats()
	}

	// now
	startTime := time.Now()
	// every logIntervalMinutes minutes we log
	logIntervalMinutes := viper.GetInt("app.logIntervalMinutes")

	// get current user
	currentUser, err := utils.GetCurrentUser()
	if err != nil {
		panic("Cannot get current user")
	}

	// Notify user about remaining time
	timeScreenLimit, err := utils.GetScreenTimeLimitMinutes(currentUser)
	if err != nil {
		logger.ZapLog.Error("Cannot get screen time limit for user ", currentUser, err)
	} else {
		totalMinutes, err := utils.GetTotalTodayTimeMinutes(currentUser)
		if err != nil {
			logger.ZapLog.Error("Cannot get today total time ", err)
		} else {
			utils.Notify(fmt.Sprintf("Hey rapoide ti tengo d'occhio, di oggi hai ancora %d minuti", timeScreenLimit-totalMinutes))
		}
	}

	// keep program running
	for {
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
				} else if totalMinutes+logIntervalMinutes > timeScreenLimit {
					utils.Notify(fmt.Sprintf("Hey rapoide stai facendo il furbo, entro %d minuti ti sloggo!", logIntervalMinutes))
				}
			}
		}

		time.Sleep(time.Duration(logIntervalMinutes) * time.Minute)
	}
}
