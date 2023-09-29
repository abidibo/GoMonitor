package core

import (
	"os"

	"github.com/abidibo/gomonitor/core/monitor"
	"github.com/abidibo/gomonitor/core/stats"
)

func Run() {
	if len(os.Args) > 1 {
		// show stats
		stats.Stats()
	} else {
		// run monitor
		monitor.Run()
	}
}
