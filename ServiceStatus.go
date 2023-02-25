package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Status int
type Action int

const (
	Inactive Status = iota
	Active
	Deactivating
)

const (
	Start Action = iota
	Stop
	Restart
)

var (
	capabilitiesMap = map[string]Status{
		"inactive":     Inactive,
		"active":       Active,
		"deactivating": Deactivating,
	}
)

var (
	capabilitiesStateMap = map[string]string{
		"inactive":     "off",
		"active":       "on",
		"deactivating": "off",
	}
)

var (
	actionMap = map[string]Action{
		"start":   Start,
		"stop":    Stop,
		"restart": Restart,
	}
)

var (
	actionMapStatus = map[string]string{
		"on":  "start",
		"off": "stop",
	}
)

type Service struct {
	name string
}

func ParseStatus(str string) (Status, bool) {
	c, ok := capabilitiesMap[strings.TrimSpace(strings.ToLower(str))]
	return c, ok
}

func ParseStatusState(str string) (string, bool) {
	c, ok := capabilitiesStateMap[strings.TrimSpace(strings.ToLower(str))]
	return c, ok
}

func (o *Service) checkStatus() (Status, bool, string, string) {
	cmd := exec.Command("systemctl", "check", o.name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			fmt.Printf("failed to run systemctl: %v", err)
			os.Exit(1)
		}
	}
	var status, error = ParseStatus(string(out))
	var statusState, _ = ParseStatusState(string(out))
	return status, error, string(out), statusState
}

func (o *Service) setStateMessage(setStateMessage string) {
	split := strings.Split(setStateMessage, "/")
	var command = strings.TrimSpace(strings.ToLower(split[len(split)-1]))
	_, ok := actionMap[command]
	if !ok {
		command, ok = actionMapStatus[command]
		if !ok {
			fmt.Printf("Service command not found: %s", command)
			return
		}
	}

	cmd := exec.Command("systemctl", command, o.name)
	cmd.CombinedOutput()
}

func (o *Service) setStatePayload(payload []byte) {
	var command = strings.TrimSpace(strings.ToLower(string(payload)))
	_, ok := actionMap[command]
	if !ok {
		command, ok = actionMapStatus[command]
		if !ok {
			fmt.Printf("Service command not found: %s", command)
			return
		}
	}

	cmd := exec.Command("systemctl", command, o.name)
	cmd.CombinedOutput()
}
