package runner

import (
	"io"
	"os/exec"
	"syscall"
	"time"
)

func run() bool {
	runnerLog("Running...")

	var cmd *exec.Cmd
	if isDelve() {
		runnerLog("... using delve")
		cmd = exec.Command(
			"dlv",
			"--listen=:2345",
			"--headless=true",
			"--accept-multiclient",
			"--api-version=2",
			"--log",
			"exec",
			buildPath(),
		)
	} else {
		cmd = exec.Command(buildPath())
	}

	cmd.Dir = workDir()

	runnerLog("Set workdir to " + workDir())

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fatal(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		fatal(err)
	}

	runnerLog("Process started PID %d", cmd.Process.Pid)

	go io.Copy(appLogWriter{}, stderr)
	go io.Copy(appLogWriter{}, stdout)

	waitChannel := make(chan bool)

	go func() {
		defer close(waitChannel)

		_ = cmd.Wait()
		runnerLog("Process exited PID %d", cmd.Process.Pid)
		waitChannel <- true
	}()

	go func() {
		<-stopChannel
		pid := cmd.Process.Pid
		runnerLog("Send sigterm to PID %d", pid)
		err = cmd.Process.Signal(syscall.SIGTERM)
		if err != nil {
			runnerLog("Failed to send sigterm to PID %d", pid)
		}

		select {
		case <-waitChannel:
		case <-time.After(time.Second * 3):
			runnerLog("Timed out waiting for process to exit PID %d", pid)
			_ = cmd.Process.Kill()
		}

		stoppedChannel <- true
	}()

	return true
}
