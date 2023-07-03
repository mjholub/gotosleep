package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/avast/retry-go/v4"
)

func main() {
	lockCmd, err := evalLockCmd()
	if err != nil {
		fmt.Println("Failed to evaluate lock command:", err)
		os.Exit(1)
	}

	// Run swaylock
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err = runCommand(lockCmd); err != nil {
			fmt.Println("Failed to run swaylock:", err)
			os.Exit(1)
		}
	}()

	wg.Wait()

	// Echo "mem" to /sys/power/state
	err = echoMemToPowerState()
	if err != nil {
		lockErr := killLocker(lockCmd)
		if lockErr != nil {
			log.Println("Failed to kill screen locker:", lockErr)
		}
		log.Fatalf("Failed to echo mem to /sys/power/state:", err)
	}
}

func runCommand(name string, arg ...string) (err error) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}

	// Poll the process to check if it's running
	checkProcess := func() error {
		if cmd.Process == nil {
			return errors.New("screen locker not running")
		}

		// Use the kill system call to check if the process is running
		err = cmd.Process.Signal(syscall.Signal(0))
		if err != nil {
			return errors.New("screen locker not running")
		}

		return nil
	}

	// Use retry-go to retry the checkProcess function with exponential backoff
	err = retry.Do(
		checkProcess,
		retry.Attempts(9),
		retry.Delay(25*time.Millisecond),
		retry.MaxDelay(4*time.Second),
		retry.DelayType(retry.BackOffDelay),
	)

	if err != nil {
		return errors.New("screen locker not running")
	}
	return nil
}

func killLocker(name string) error {
	// Get the process ID of the screen locker
	out := exec.Command("pidof", name)
	var outb, errb bytes.Buffer
	out.Stdout = &outb
	out.Stderr = &errb
	if err := out.Run(); err != nil {
		return errors.New("failed to get screen locker PID")
	}
	pidStr := strings.TrimSpace(outb.String())
	if pid, err := strconv.Atoi(pidStr); err != nil {
		return errors.New("failed to convert screen locker PID to int")
	} else if err := syscall.Kill(int(pid), syscall.SIGTERM); err != nil {
		return errors.New("failed to kill screen locker")
	}
	return nil
}
