package main

import (
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/alt-dima/go-entrypoint/utils"
)

var (
	version string = "unspecified"
)

func main() {
	var logger = utils.Logger
	slog.SetDefault(logger)

	if len(os.Args) == 1 {
		logger.Error("Entrypoint nothing to be executed")
		os.Exit(1)
	} else if os.Args[1] == "" {
		logger.Error("Entrypoint empty command")
		os.Exit(1)
	}
	cmdToExec := os.Args[1]

	argsToExec := []string{}
	for _, value := range os.Args[2:] {
		if value != "" {
			argsToExec = append(argsToExec, value)
		}
	}

	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	logger.Debug("Entrypoint version " + version)

	logger.Debug("Entrypoint starting child: " + cmdToExec + " " + strings.Join(argsToExec, " "))
	cmd := exec.Command(cmdToExec, argsToExec...)
	cmd.Env = utils.GenerateChildEnvs()
	// pipe the commands output to the applications
	// standard output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		logger.Error("Entrypoint failed to start child: " + err.Error())
		os.Exit(1)
	}

	var shutdownStartTime time.Time

	go func() {
		sig := <-sigs
		logger.Debug("Entrypoint got signal " + sig.String())
		shutdownStartTime = time.Now()
		cmd.Process.Signal(sig)
	}()
	//log.Println("awaiting signal")

	err = cmd.Wait()
	exitCodeFinal := 0
	if err != nil && cmd.ProcessState.ExitCode() < 0 {
		exitCodeFinal = 1
		logger.Warn("Entrypoint child failed: " + err.Error())
	} else if cmd.ProcessState.ExitCode() == 143 {
		exitCodeFinal = 0
	} else {
		exitCodeFinal = cmd.ProcessState.ExitCode()
	}

	// Could be used to stop sidecars such envoy proxy by sending specific HTTP-request
	//utils.StopExtSvcs()

	if !shutdownStartTime.IsZero() {
		shutdownElapsedDuration := time.Since(shutdownStartTime)
		if shutdownElapsedDuration.Seconds() > 85 {
			logger.Warn("Entrypoint slow child termination took " + shutdownElapsedDuration.String())
		} else if shutdownElapsedDuration.Milliseconds() < 50 {
			logger.Info("Entrypoint fast child termination took " + shutdownElapsedDuration.String())
		}
	}
	if exitCodeFinal == 0 {
		logger.Debug("Entrypoint exiting with code " + strconv.Itoa(exitCodeFinal))
	} else {
		logger.Warn("Entrypoint exiting with code " + strconv.Itoa(exitCodeFinal))
	}
	os.Exit(exitCodeFinal)
}
