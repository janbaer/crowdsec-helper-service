package csclirunner

import (
	"log/slog"
	"os"
	"os/exec"
)

// const CSLI_BINARY_PATH = "/usr/bin/cscli"
const CSLI_BINARY_PATH = "echo"

var logger *slog.Logger

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func DeleteDecision(ipAddress string) error {
	// logger.Info("Deleting decision for IP", "ipAddress", ipAddress)
	// command = "/usr/bin/cscli decisions delete -i="..query.ip
	return executeCommand("decisions", "delete", "i="+ipAddress)
}

func CreateDecision(ipAddress string, decisionType string, duration string) error {
	logger.Info(
		"Creating decision for IP: ",
		"ipAddress", ipAddress,
		"type", decisionType,
		"duration", duration,
	)
	return executeCommand(
		"decisions",
		"add",
		"i="+ipAddress,
		"type="+decisionType,
		"duration="+duration,
	)
}

func executeCommand(command string, subCommand string, args ...string) error {
	execArgs := append([]string{command, subCommand}, args...)
	cmd := exec.Command(CSLI_BINARY_PATH, execArgs...)
	if err := cmd.Run(); err != nil {
		logger.Error("Error running cscli with args", "args", args, "err", err)
		return err
	}
	logger.Info("Executed command with args", "binary", CSLI_BINARY_PATH, "args", execArgs)
	return nil
}
