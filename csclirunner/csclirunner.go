package csclirunner

import (
	"log/slog"
	"os"
	"os/exec"
)

const CSLI_BINARY_PATH = "/usr/bin/cscli"

var logger *slog.Logger

func init() {
	// logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func DeleteDecision(ipAddress string) error {
	logger.Info("Deleting decision for IP", "ipAddress", ipAddress)
	return executeCommand("decisions", "delete", "--ip="+ipAddress)
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
		"--ip="+ipAddress,
		"--type="+decisionType,
		"--duration="+duration,
	)
}

func executeCommand(command string, subCommand string, args ...string) error {
	execArgs := append([]string{command, subCommand}, args...)
	cmd := exec.Command(CSLI_BINARY_PATH, execArgs...)
	if err := cmd.Run(); err != nil {
		logger.Error("Error running cscli with args", "args", execArgs, "err", err)
		return err
	}
	logger.Info("Executed command with args", "binary", CSLI_BINARY_PATH, "args", execArgs)
	return nil
}
