package monitor

import (
	"github.com/abidibo/gomonitor/core/utils"
)

func Run() {
	if utils.IsRunningAsRoot() {
		RunAsRoot()
	} else {
		RunNonRoot()
	}
}
