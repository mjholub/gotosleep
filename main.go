package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
)

func main() {
	lockCmd, err := evalLockCmd()
	if err != nil {
		fmt.Println("Failed to evaluate lock command:", err)
		os.Exit(1)
	}

	// Run swaylock
	if err = runCommand(lockCmd); err != nil {
		fmt.Println("Failed to run swaylock:", err)
		os.Exit(1)
	}

	// Check if the current user is root
	// FIXME: add polkit/elogind/dbus based authentication
	if !isRoot() {
		fmt.Println("You must run this command as root.")
		os.Exit(1)
	}
	// Echo "mem" to /sys/power/state
	err = echoMemToPowerState()
	if err != nil {
		fmt.Println("Failed to echo mem to /sys/power/state:", err)
		os.Exit(1)
	}
}

func isRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Failed to get current user:", err)
		os.Exit(1)
	}
	return currentUser.Uid == "0"
}

func runCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func echoMemToPowerState() error {
	// Open the file
	file, err := os.OpenFile("/sys/power/state", os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write "mem" to the file
	_, err = file.WriteString("mem")
	return err
}
