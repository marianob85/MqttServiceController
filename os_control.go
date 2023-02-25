package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func osControlPayload(payload []byte) {
	if strings.TrimSpace(strings.ToLower(string(payload))) == "reboot" {
		fmt.Println("System reboot...")
		cmd := exec.Command("reboot", "0")
		cmd.CombinedOutput()
	}
}

func osControlMessage(topic string) {
	split := strings.Split(topic, "/")
	if strings.TrimSpace(strings.ToLower(split[len(split)-1])) == "reboot" {
		fmt.Println("System reboot...")
		cmd := exec.Command("reboot", "0")
		cmd.CombinedOutput()
	}
}
