package monitor

import (
	"fmt"
	"image/color"
	"sort"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/abidibo/gomonitor/core"
	"github.com/abidibo/gomonitor/core/utils"
	"github.com/abidibo/gomonitor/logger"
	"github.com/spf13/viper"
)

var a fyne.App
var w fyne.Window

// entry point
func RunNonRoot() {
	// spawn a thread to monitor
	go notificationsThread()
	createMainWindow()
}

func createMainWindow() {
	a = app.NewWithID("net.abidibo.gomonitor")
	w = a.NewWindow("GoMonitor")

	if desk, ok := a.(desktop.App); ok {
		m := fyne.NewMenu(
			"GoMonitor",
			fyne.NewMenuItem("Stats", func() {
				w.Show()
			}),
			fyne.NewMenuItem("Quit", func() {
				w.Hide()
			}),
		)
		desk.SetSystemTrayMenu(m)
		icon := fyne.NewStaticResource("goMonitorIcon", core.IconData)
		desk.SetSystemTrayIcon(icon)
	}

	updateWindowContent()
	go keepUpdatingWindow()

	w.SetCloseIntercept(func() {
		w.Hide()
	})
	w.Resize(fyne.NewSize(600, 600))
	w.ShowAndRun()
}

func keepUpdatingWindow() {
	for range time.Tick(time.Second * 60) {
		updateWindowContent()
	}
}

func updateWindowContent() {
	// now label
	// #ff9900 in rgba: 255, 153, 0
	nowLabel := canvas.NewText(fmt.Sprintf("Last update: %s", time.Now().Format("15:04:05")), color.RGBA{255, 153, 0, 255})
	// today usage time
	user, _ := utils.GetCurrentUser()
	screenTimeConfiguration, err := utils.GetScreenTimeConfiguration(user)
	totalTodayMinutes, _ := utils.GetTotalTodayTimeMinutes(user)
	totalTodayTimeLabel := canvas.NewText(fmt.Sprintf("Total minutes today: %d/%d", totalTodayMinutes, screenTimeConfiguration.ScreenLimitMin), color.RGBA{255, 153, 0, 255})
	totalTodayTimeLabel.TextStyle.Bold = true

	screentTimeWindowLabel := canvas.NewText(fmt.Sprintf("Allowed time: %s/%s", screenTimeConfiguration.TimeWindowStart, screenTimeConfiguration.TimeWindowStop), color.RGBA{255, 153, 0, 255})

	var data [][]string = nil
	processTotalMap := make(map[string]int)
	processesBodyMap := make(map[string][]string)
	dims := []int{0, 0}
	date := time.Now().Format("2006-01-02")
	processes, err := utils.GetAllDateProcesses(user, date, 20)

	if err != nil {
		fmt.Println(err)
	} else {
		dims[0] = len(processes)
		dims[1] = 2
		data = append(data, []string{"Process", "Time (min)"})
		for _, p := range processes {
			total, err := utils.GetTotalProcessTimeMinutes(user, p, date)
			if err != nil {
				fmt.Println(err)
			} else {
				processTotalMap[p] = total
				processesBodyMap[p] = []string{fmt.Sprintf("%s", p), fmt.Sprintf("%d", total)}
			}
		}

		// sort by total time
		keys := make([]string, 0, len(processTotalMap))
		for k := range processTotalMap {
			keys = append(keys, k)
		}
		sort.SliceStable(keys, func(i, j int) bool {
			return processTotalMap[keys[i]] > processTotalMap[keys[j]]
		})
		for _, k := range keys {
			data = append(data, processesBodyMap[k])
		}
	}

	list := widget.NewTable(
		func() (int, int) {
			return dims[0], dims[1]
		},
		func() fyne.CanvasObject {
			item := widget.NewLabel("wide content")
			return item
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[i.Row][i.Col])
			if i.Row == 0 {
				o.(*widget.Label).TextStyle.Bold = true
				o.(*widget.Label).TextStyle.Italic = false
				o.(*widget.Label).Alignment = fyne.TextAlignLeading
			} else if i.Col == 0 {
				o.(*widget.Label).TextStyle.Italic = true
				o.(*widget.Label).TextStyle.Bold = false
				o.(*widget.Label).Alignment = fyne.TextAlignLeading
			} else if i.Col == 1 {
				o.(*widget.Label).TextStyle.Italic = false
				o.(*widget.Label).TextStyle.Bold = false
				o.(*widget.Label).Alignment = fyne.TextAlignTrailing
			} else {
				o.(*widget.Label).Alignment = fyne.TextAlignLeading
				o.(*widget.Label).TextStyle.Italic = false
				o.(*widget.Label).TextStyle.Bold = false
			}
		})
	list.SetColumnWidth(0, 480)

	header := container.New(
		layout.NewVBoxLayout(),
		container.New(layout.NewHBoxLayout(), layout.NewSpacer(), totalTodayTimeLabel, layout.NewSpacer()),
		container.New(layout.NewHBoxLayout(), layout.NewSpacer(), screentTimeWindowLabel, layout.NewSpacer()),
	)
	footer := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), nowLabel, layout.NewSpacer())
	table := container.New(layout.NewStackLayout(), list)
	content := container.New(layout.NewBorderLayout(header, footer, nil, nil), header, footer, table)
	w.SetContent(content)
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
	screenTimeConfiguration, err := utils.GetScreenTimeConfiguration(currentUser)
	if err != nil {
		logger.ZapLog.Error("Cannot get screen time limit for user ", currentUser, err)
	} else {
		totalMinutes, err := utils.GetTotalTodayTimeMinutes(currentUser)
		if err != nil {
			logger.ZapLog.Error("Cannot get today total time ", err)
		} else {
			utils.Notify(fmt.Sprintf("Hey rapoide ti tengo d'occhio, di oggi hai ancora %d minuti", screenTimeConfiguration.ScreenLimitMin-totalMinutes))
		}
	}

	timeWindowAlertSent := false
	for {
		// if taken from api, we need to refresh it
		screenTimeConfiguration, err := utils.GetScreenTimeConfiguration(currentUser)
		// get sum of partial_time_minutes for current day and user
		totalMinutes, err := utils.GetTotalTodayTimeMinutes(currentUser)
		if err == nil {
			// check if total time usage has exceeded the limit
			if screenTimeConfiguration.ScreenLimitMin > 0 {
				if (totalMinutes < screenTimeConfiguration.ScreenLimitMin/2) && (totalMinutes+logIntervalMinutes > screenTimeConfiguration.ScreenLimitMin/2) {
					utils.Notify("Hey rapoide sei più o meno a metà del tuo tempo a disposizione, fai cisti!")
				}
				if totalMinutes+logIntervalMinutes > screenTimeConfiguration.ScreenLimitMin {
					utils.Notify(fmt.Sprintf("Hey rapoide stai facendo il furbo, entro %d minuti ti sloggo!", logIntervalMinutes))
				}
			}
		}

		// check time window
		now := time.Now()
		nowInt, _ := strconv.Atoi(strings.Replace(fmt.Sprintf("%02d:%02d", now.Hour(), now.Minute()), ":", "", 1))
		startInt, err := strconv.Atoi(strings.Replace(screenTimeConfiguration.TimeWindowStart, ":", "", 1))
		if err == nil {
			stopInt, err := strconv.Atoi(strings.Replace(screenTimeConfiguration.TimeWindowStop, ":", "", 1))
			if err == nil {
				if startInt != stopInt && stopInt-nowInt < 10 && !timeWindowAlertSent {
					timeWindowAlertSent = true
					utils.Notify(fmt.Sprintf("Hey rapoide guarda che mancano meno di 10 minuti, ti conviene salvare il salvabile!"))
				}
			}
		}

		time.Sleep(time.Duration(logIntervalMinutes) * time.Minute)
	}

}
