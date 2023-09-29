package utils

import (
	"os"
	"os/exec"
	"os/user"
	"strings"
	"syscall"

	"github.com/abidibo/gomonitor/logger"
)

func IsRunningAsRoot() bool {
	return os.Geteuid() == 0
}

func GetRunningUsername() (string, error) {
	pUser, err := user.Current()
	if err != nil {
		return "", nil
	}
	return pUser.Username, nil
}

func GetCurrentUser() (string, error) {
	cmd := exec.Command("who")
	out, err := cmd.Output()
	if err != nil {
		logger.ZapLog.Error("Cannot get current user")
		return "", err
	} else {
		logger.ZapLog.Info("Current user ", string(out))
		parts := strings.Split(string(out), " ")
		return parts[0], nil
	}
}

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
