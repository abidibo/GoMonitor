package stats

import (
	"fmt"
	"os"
	"time"

	"github.com/abidibo/gomonitor/core/utils"
)

func Stats() {

	var user string
	var date string

	if len(os.Args) > 2 && os.Args[1] == "stats" {
		if len(os.Args) < 4 {
			currentTime := time.Now()
			user = os.Args[2]
			date = currentTime.Format("2006-01-02")
		} else {
			user = os.Args[2]
			date = os.Args[3]
		}
	} else if len(os.Args) > 1 {
		utils.Help()
		return
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
