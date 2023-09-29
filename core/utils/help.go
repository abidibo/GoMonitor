package utils

import "fmt"

func Help() {
	fmt.Println("Usage to run the monitor (as root): gomonitor")
	fmt.Println("Usage to view statistics (as root): gomonitor stats [user] [date?]")
	fmt.Println("Usage to view statistics (not root): gomonitor stats [date?]")
}
