package monitor

import (
	"bytes"
	"fmt"
	"image/png"
	"time"

	"github.com/abidibo/gomonitor/core"
	"github.com/abidibo/gomonitor/core/utils"
	"github.com/abidibo/gomonitor/logger"
	"github.com/gen2brain/iup-go/iup"
	"github.com/spf13/viper"
)

// entry point
func RunNonRoot() {
	// spawn a thread to monitor
	go notificationsThread()

	// main application window
	iup.Open()
	// closing the application, do not stop the main loop!
	iup.SetGlobal("LOCKLOOP", "YES")
	// create window
	createMainWindow()
	// main loop
	iup.MainLoop()
}

func createMainWindow() {
	img, _ := png.Decode(bytes.NewReader(core.IconData))
	iup.ImageFromImage(img).SetHandle("goMonitorIcon")

	// today usage time
	user, _ := utils.GetCurrentUser()
	timeScreenLimit, _ := utils.GetScreenTimeLimitMinutes(user)
	totalTodayMinutes, _ := utils.GetTotalTodayTimeMinutes(user)
	totalTodayTimeLabel := iup.Label("Total minutes today").SetAttributes(`EXPAND=YES, ALIGNMENT=ACENTER`)
	totalTodayTimeValue := iup.Label(fmt.Sprintf("%d/%d", totalTodayMinutes, timeScreenLimit)).SetAttributes(`EXPAND=YES, ALIGNMENT=ACENTER, PADDING=10x10`)
	// btnRefresh := iup.Button("Refresh").SetAttributes(`EXPAND=YES, ALIGNMENT=ACENTER`)

	// iup.SetCallback(btnRefresh, "ACTION", iup.ActionFunc(btnRefreshCb))

	iup.SetHandle("totalTodayTimeValue", totalTodayTimeValue)

	var columns []iup.Ihandle
	columns = append(columns, iup.Label("Process").SetAttributes("FONTSTYLE=Bold"))
	columns = append(columns, iup.Label("Time (min)").SetAttributes("FONTSTYLE=Bold"))
	date := time.Now().Format("2006-01-02")
	processes, err := utils.GetAllDateProcesses(user, date, 20)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, p := range processes {
			total, err := utils.GetTotalProcessTimeMinutes(user, p, date)
			if err != nil {
				fmt.Println(err)
			} else {
				columns = append(columns, iup.Label(fmt.Sprintf("%s", p)))
				columns = append(columns, iup.Label(fmt.Sprintf("%d", total)))
			}
		}
	}

	fr := iup.Frame(
		iup.GridBox(
			columns...,
		).SetAttributes(map[string]string{
			"ALIGNMENTCOL1":  "ARIGHT",
			"HOMOGENEOUSLIN": "Yes",
			"HOMOGENEOUSCOL": "Yes",
			"NUMDIV":         "2",
			"MARGIN":         "10x10",
			"GAPLIN":         "5",
			"GAPCOL":         "15",
			"SIZE":           fmt.Sprintf("%dx%d", 230, len(columns)*5),
		}),
	)

	win := iup.Dialog(
		iup.Vbox(
			iup.Vbox(
				totalTodayTimeLabel,
				totalTodayTimeValue,
			),
			fr,
		).SetAttributes(`MARGIN=20x20`),
	).SetAttributes(map[string]string{
		"TITLE":     "GoMonitor",
		"TRAY":      "YES",
		"TRAYTIP":   "The best monitor app in the world",
		"TRAYIMAGE": "goMonitorIcon",
		"ICON":      "goMonitorIcon",
		"SIZE":      "260x300",
		"RESIZE":    "YES",
	}).SetCallback("TRAYCLICK_CB", iup.TrayClickFunc(trayClickCb)).SetHandle("win")
	// trick to open the main window, dhow tray icon and close it
	iup.Show(win)
	iup.Hide(win)
}

func btnRefreshCb(ih iup.Ihandle) int {
	updateMainWindow()
	return iup.DEFAULT
}
func updateMainWindow() {
	user, _ := utils.GetCurrentUser()
	timeScreenLimit, _ := utils.GetScreenTimeLimitMinutes(user)
	totalTodayMinutes, _ := utils.GetTotalTodayTimeMinutes(user)
	totalTodayTimeValue := iup.GetHandle("totalTodayTimeValue")
	totalTodayTimeValue.SetAttribute("TITLE", fmt.Sprintf("%d/%d", totalTodayMinutes, timeScreenLimit))
}

func trayClickCb(ih iup.Ihandle, but, pressed, dclick int) int {
	if but == 1 && pressed > 0 {
		updateMainWindow()
		iup.Show(iup.GetHandle("win"))
	}
	return iup.DEFAULT
}

// notify user at the beginning, half time and end time reached
func notificationsThread() {
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

		time.Sleep(time.Duration(logIntervalMinutes) * time.Second)
	}

}
