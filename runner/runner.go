package runner

import (
	"io"
	"os/exec"
	"syscall"
	"time"
)

func run() bool {
	runnerLog("Running...")

	cmd := exec.Command(buildPath())
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

	go io.Copy(appLogWriter{}, stderr)
	go io.Copy(appLogWriter{}, stdout)

	go func() {
		<-stopChannel
		pid := cmd.Process.Pid
		runnerLog("Send sigterm to PID %d", pid)
		err = cmd.Process.Signal(syscall.SIGTERM)
		if err != nil {
			runnerLog("Failed to send sigterm to PID %d", pid)
		}

		waitChannel := make(chan bool)

		go func() {
			defer close(waitChannel)

			_ = cmd.Wait()
			waitChannel <- true
		}()

		select {
		case <-waitChannel:
			runnerLog("Process exited PID %d", pid)
		case <-time.After(time.Second * 3):
			runnerLog("Timed out waiting for process to exit PID %d", pid)
			_ = cmd.Process.Kill()
		}

		stoppedChannel <- true
	}()

	return true
}
