package main

import (
	"os/exec"
	"strings"
)

func osControl(payload []byte) {
	if strings.TrimSpace(strings.ToLower(string(payload))) == "restart" {
		cmd := exec.Command("reboot", "0")
		cmd.CombinedOutput()
	}
}
