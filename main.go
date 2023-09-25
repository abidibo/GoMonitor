package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"os/user"
	"sort"
	"syscall"
	"time"

	"github.com/TheCreeper/go-notify"
	"github.com/abidibo/gomonitor/db"
	"github.com/abidibo/gomonitor/logger"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/spf13/viper"
)

func init() {
	// Read settings
	viper.SetConfigFile(fmt.Sprintf("/etc/gomonitor.json"))
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("Error reading settings file, %s", err))
	}

	// Init logger
	logger.InitLogger()
	logger.ZapLog.Info("Starting gomonitor application")

	// Ensure app home directory exists
	homePath := viper.GetString("app.homePath")

	err := os.MkdirAll(homePath, os.ModePerm)
	if err != nil {
		panic("Error creating the application home directory")
	} else {
		logger.ZapLog.Info("Application home directory ok ", homePath)
	}
}

type ProcessLog struct {
	Name          string
	CpuPercent    float64
	MemoryPercent float32
}

func main() {
	// Init database
	db.InitDatabase()

	// time control stuff
	startTime := time.Now()
	logIntervalMinutes := viper.GetInt("app.logIntervalMinutes")

	// get current user
	logger.ZapLog.Debug("Retrieving current user")
	currentUser, err := user.Current()
	if err != nil {
		logger.ZapLog.Error("Cannot get current user")
	} else {
		logger.ZapLog.Debug("Current user ", currentUser.Username)
	}

	// get screen time limit
	limits := make(map[string]int)
	err = viper.UnmarshalKey("app.screenTimeLimitMinutes", &limits)
	timeScreenLimit, okTimeScreenLimit := limits[currentUser.Username]

	// get remaining time
	stm, err := db.DB().C.Prepare("SELECT SUM(total_time_minutes) FROM log WHERE user = ? AND date(timestamp) = date('now')")
	if err != nil {
		logger.ZapLog.Error("Cannot get total time", err)
	} else if okTimeScreenLimit {
		var totalMinutes int
		stm.QueryRow(currentUser.Username).Scan(&totalMinutes)
		ntf := notify.NewNotification("GoMonitor", fmt.Sprintf("Hey rapoide ti tengo d'occhio, di oggi hai ancora %d minuti", timeScreenLimit-totalMinutes))
		if _, err := ntf.Show(); err != nil {
			logger.ZapLog.Error("Cannot show notification ", err)
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
		result, err := db.DB().C.Exec("INSERT INTO log (user, total_time_minutes) VALUES (?, ?)", currentUser.Username, elapsed)
		startTime = time.Now() // reset start time
		if err != nil {
			logger.ZapLog.Error("Cannot insert log into DB", err)
		} else {
			logId, _ := result.LastInsertId()
			for _, p := range logProcesses {
				db.DB().C.Exec("INSERT INTO log_process (log_id, name, cpu_percent, memory_percent) VALUES (?, ?, ?, ?)", logId, p.Name, p.CpuPercent, p.MemoryPercent)
			}
		}

		// get sum of total_time_minutes for current day and user
		stm, err := db.DB().C.Prepare("SELECT SUM(total_time_minutes) FROM log WHERE user = ? AND date(timestamp) = date('now')")
		if err != nil {
			logger.ZapLog.Error("Cannot get total time", err)
		} else {
			var totalMinutes int
			stm.QueryRow(currentUser.Username).Scan(&totalMinutes)
			logger.ZapLog.Info("Total time for user ", currentUser.Username, " ", totalMinutes, " minutes")

			// check if total time usage has exceeded the limit
			if okTimeScreenLimit {
				if totalMinutes > timeScreenLimit {
					logger.ZapLog.Info("Limit exceeded for user ", currentUser.Username, " ", totalMinutes, " minutes")
					// logout
					c := exec.Command("pkill", "-KILL", "-u", currentUser.Username)
					err = c.Run()
					if err != nil {
						logger.ZapLog.Error("Cannot logout user ", currentUser.Username, err)
						// shutdown pc
						err = syscall.Reboot(syscall.LINUX_REBOOT_CMD_HALT)
						if err != nil {
							logger.ZapLog.Error("Cannot halt pc ", err)

							err = syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
							if err != nil {
								logger.ZapLog.Error("Cannot shutdown pc ", err)
							}
						}
					}
				} else if totalMinutes+logIntervalMinutes > timeScreenLimit {
					ntf := notify.NewNotification("GoMonitor", fmt.Sprintf("Hey rapoide stai facendo il furbo, entro %d minuti ti sloggo!", logIntervalMinutes))
					if _, err := ntf.Show(); err != nil {
						logger.ZapLog.Error("Cannot show notification ", err)
					}
				}
			}
		}

		time.Sleep(time.Duration(logIntervalMinutes) * time.Minute)
	}
}
