package utils

import (
	"github.com/TheCreeper/go-notify"
	"github.com/abidibo/gomonitor/db"
	"github.com/abidibo/gomonitor/logger"
)

func GetTodayTimeMinutes(user string) (int, error) {
	stm, err := db.DB().C.Prepare("SELECT SUM(total_time_minutes) FROM log WHERE user = ? AND date(timestamp) = date('now')")
	if err != nil {
		return 0, err
	} else {
		var totalMinutes int
		stm.QueryRow(user).Scan(&totalMinutes)
		return totalMinutes, nil
	}
}

func InsertLog(user string, totalMinutes int) (int64, error) {
	result, err := db.DB().C.Exec("INSERT INTO log (user, total_time_minutes) VALUES (?, ?)", user, totalMinutes)

	if err != nil {
		logger.ZapLog.Error("Cannot insert log into DB", err)
		return 0, err
	} else {
		logId, _ := result.LastInsertId()
		return logId, nil
	}
}

func InsertProcessLog(logId int64, name string, cpuPercent float64, memoryPercent float32) error {
	_, err := db.DB().C.Exec("INSERT INTO log_process (log_id, name, cpu_percent, memory_percent) VALUES (?, ?, ?, ?)", logId, name, cpuPercent, memoryPercent)
	if err != nil {
		logger.ZapLog.Error("Cannot insert process log into DB", err)
	}

	return err
}

func GetTotalTodayTimeMinutes(user string) (int, error) {
	stm, err := db.DB().C.Prepare("SELECT SUM(total_time_minutes) FROM log WHERE user = ? AND date(timestamp) = date('now')")
	if err != nil {
		logger.ZapLog.Error("Cannot get total day time", err)
		return 0, err
	} else {
		var totalMinutes int
		stm.QueryRow(user).Scan(&totalMinutes)
		return totalMinutes, nil
	}
}

func Notify(text string) {
	ntf := notify.NewNotification("GoMonitor", text)
	if _, err := ntf.Show(); err != nil {
		logger.ZapLog.Error("Cannot show notification ", err)
	}
}
