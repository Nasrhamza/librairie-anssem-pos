//go:build !windows

package main

import (
	"fmt"
	"os/exec"
	"runtime"
)

func runDesktopWindow(url string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "darwin" {
		cmd = exec.Command("open", url)
	} else {
		cmd = exec.Command("xdg-open", url)
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	return fmt.Errorf("la fenêtre desktop native est disponible uniquement sous Windows")
}

func showDesktopError(_, message string) { fmt.Println(message) }
