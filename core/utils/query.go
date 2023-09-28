package utils

import (
	"github.com/TheCreeper/go-notify"
	"github.com/abidibo/gomonitor/db"
	"github.com/abidibo/gomonitor/logger"
)

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

func GetTotalDateTimeMinutes(user string, date string) (int, error) {
	stm, err := db.DB().C.Prepare("SELECT SUM(total_time_minutes) FROM log WHERE user = ? AND date(timestamp) = ?")
	if err != nil {
		logger.ZapLog.Error("Cannot get total day time", err)
		return 0, err
	} else {
		var totalMinutes int
		stm.QueryRow(user, date).Scan(&totalMinutes)
		return totalMinutes, nil
	}
}

func GetTotalProcessTimeMinutes(user string, processName string, date string) (int, error) {
	// all processes
	// SELECT DISTINCT name, SUM(total_time_minutes) as tot FROM (SELECT p.id, p.log_id, p.name, l.total_time_minutes from log_process as p, log as l WHERE l.user = 'abidibo' AND date(l.timestamp) = '2023-09-25' AND p.log_id = l.id GROUP BY p.log_id, p.name) GROUP BY name ORDER BY tot DESC;
	stm, err := db.DB().C.Prepare("SELECT SUM(total_time_minutes) FROM (SELECT p.log_id, p.name, l.total_time_minutes from log_process as p, log as l WHERE p.name = ? and l.user = ? AND date(l.timestamp) = ? AND p.log_id = l.id GROUP BY p.log_id, p.name)")
	if err != nil {
		logger.ZapLog.Error("Cannot get total time per process", err)
		return 0, err
	} else {
		var totalMinutes int
		stm.QueryRow(processName, user, date).Scan(&totalMinutes)
		return totalMinutes, nil
	}
}

func GetAllDateProcesses(user string, date string, limit int) ([]string, error) {
	stm, err := db.DB().C.Prepare("SELECT DISTINCT p.name FROM log_process AS p, log as l WHERE l.user = ? AND date(l.timestamp) = ? ORDER BY p.cpu_percent DESC, p.memory_percent DESC, p.name ASC LIMIT 0, ?;")
	if err != nil {
		logger.ZapLog.Error("Cannot get all date processes", err)
		return nil, err
	} else {
		var processes []string
		rows, err := stm.Query(user, date, limit)
		if err != nil {
			logger.ZapLog.Error("Cannot get all date processes", err)
			return nil, err
		} else {
			for rows.Next() {
				var process string
				rows.Scan(&process)
				processes = append(processes, process)
			}
			return processes, nil
		}
	}
}

func Notify(text string) {
	ntf := notify.NewNotification("GoMonitor", text)
	if _, err := ntf.Show(); err != nil {
		logger.ZapLog.Error("Cannot show notification ", err)
	}
}
