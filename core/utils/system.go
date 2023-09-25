package utils

import (
	"os/exec"
	"syscall"

	"github.com/abidibo/gomonitor/logger"
)

func LogoutUser(user string) error {
	cmd := exec.Command("pkill", "-KILL", "-u", user)
	err := cmd.Run()
	if err != nil {
		logger.ZapLog.Error("Cannot logout user ", user, err)
		return err
	}
	return nil
}

func Shutdown() error {
	err := syscall.Reboot(syscall.LINUX_REBOOT_CMD_HALT)
	if err != nil {
		logger.ZapLog.Error("Cannot halt pc ", err)

		err = syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
		if err != nil {
			logger.ZapLog.Error("Cannot shutdown pc ", err)
			return err
		}

		return nil
	}

	return nil
}
