package monitor

import (
	"fmt"
	"time"

	"github.com/abidibo/gomonitor/core/utils"
	"github.com/abidibo/gomonitor/logger"
	"github.com/spf13/viper"
)

func RunNonRoot() {
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

	for {
		// get sum of partial_time_minutes for current day and user
		totalMinutes, err := utils.GetTotalTodayTimeMinutes(currentUser)
		if err == nil {
			// check if total time usage has exceeded the limit
			if timeScreenLimit > 0 {
				if (totalMinutes < timeScreenLimit/2) && (totalMinutes+logIntervalMinutes > timeScreenLimit/2) {
					utils.Notify("Hey rapoide sei più o meno a metà del tuo tempo a disposizione, fai cisti!")
				}
				if totalMinutes+logIntervalMinutes > timeScreenLimit {
					utils.Notify(fmt.Sprintf("Hey rapoide stai facendo il furbo, entro %d minuti ti sloggo!", logIntervalMinutes))
				}
			}
		}

		time.Sleep(time.Duration(logIntervalMinutes) * time.Minute)
	}
}
