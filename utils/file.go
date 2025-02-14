package utils

import (
	"os"
	"os/exec"
)

func IsExecutable(command string, args ...string) bool {
	proc := exec.Command(command, args...)
	return proc.Run() == nil
}

func IsExistsFile(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
