package stats

import (
	"fmt"

	"github.com/abidibo/gomonitor/core/utils"
)

func Run(user string, date string) {
	var err error

	if !utils.IsRunningAsRoot() {
		user, err = utils.GetRunningUsername()
		if err != nil {
			panic("Cannot get current user")
		}
	}

	fmt.Println("================================================")
	fmt.Println("Stats for ", user, " ", date)
	fmt.Println("================================================")
	total, err := utils.GetTotalDateTimeMinutes(user, date)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Total time ", fmt.Sprintf("%d", total), " minutes")
		fmt.Println()
	}

	processes, err := utils.GetAllDateProcesses(user, date, 20)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, p := range processes {
			total, err := utils.GetTotalProcessTimeMinutes(user, p, date)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(fmt.Sprintf("%-25s", p), fmt.Sprintf("%4d", total), " minutes")
			}
		}
	}

}
