package core

import (
	"fmt"
	"math"
	"os/user"
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

func Run() {
	// now
	startTime := time.Now()
	// every logIntervalMinutes minutes we log
	logIntervalMinutes := viper.GetInt("app.logIntervalMinutes")

	// get current user
	logger.ZapLog.Debug("Retrieving current user")
	currentUser, err := user.Current()
	if err != nil {
		logger.ZapLog.Error("Cannot get current user")
	} else {
		logger.ZapLog.Info("Current user ", currentUser.Username)
	}

	// Notify user about remaining time
	timeScreenLimit, err := utils.GetScreenTimeLimitMinutes(currentUser.Username)
	if err != nil {
		logger.ZapLog.Error("Cannot get screen time limit for user ", currentUser.Username, err)
	} else {
		totalMinutes, err := utils.GetTodayTimeMinutes(currentUser.Username)
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
			if u == currentUser.Username && (cpuUsage > 0 || memoryUsage > 0) {
				// add process to list if it belongs to current user and is not in the list
				log := ProcessLog{Name: name, CpuPercent: cpuUsage, MemoryPercent: memoryUsage}
				logProcesses = append(logProcesses, log)
			}
		}
		sort.Slice(logProcesses, func(i, j int) bool {
			return logProcesses[i].CpuPercent > logProcesses[j].CpuPercent
		})

		// insert log into DB
		elapsed := int(math.Ceil(time.Since(startTime).Minutes()))
		startTime = time.Now() // reset start time
		logId, err := utils.InsertLog(currentUser.Username, elapsed)
		if err == nil {
			for _, p := range logProcesses {
				utils.InsertProcessLog(logId, p.Name, p.CpuPercent, p.MemoryPercent)
			}
		}

		// get sum of total_time_minutes for current day and user
		totalMinutes, err := utils.GetTotalTodayTimeMinutes(currentUser.Username)
		if err == nil {
			// check if total time usage has exceeded the limit
			if timeScreenLimit > 0 {
				if totalMinutes > timeScreenLimit {
					logger.ZapLog.Info("Limit exceeded for user ", currentUser.Username, " ", totalMinutes, " minutes")
					// logout user
					err = utils.LogoutUser(currentUser.Username)
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

func Stats(date string) {
	fmt.Println("================================================")
	fmt.Println("Stats for ", date)
	fmt.Println("================================================")
	total, err := utils.GetTotalDateTimeMinutes("abidibo", date)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Total time ", fmt.Sprintf("%d", total), " minutes\n")
	}

	processes, err := utils.GetAllDateProcesses("abidibo", date, 20)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, p := range processes {
			total, err := utils.GetTotalProcessTimeMinutes("abidibo", p, date)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(fmt.Sprintf("%-25s", p), fmt.Sprintf("%4d", total), " minutes")
			}
		}
	}

}
