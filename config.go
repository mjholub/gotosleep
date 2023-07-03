package main

import (
	"fmt"
	"os"

	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

func checkConfFileExists() (bool, error) {
	// check if config file exists, if not, create it and write default config
	if _, err := os.Stat("/etc/gotosleep/gts.yaml"); os.IsNotExist(err) {
		// Create the config file
		configFile, err := os.Create("/etc/gotosleep/gts.yaml")
		if err != nil {
			return false, err
		}
		defer configFile.Close()

		// Write the default config to the file
		_, err = configFile.WriteString("lock_command: swaylock")
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func evalLockCmd() (lockCmd string, err error) {
	if confFileExists, err := checkConfFileExists(); !confFileExists || err != nil {
		return "", err
	}
	// Open the config file
	configFile, err := os.Open("/etc/gotosleep/gts.yaml")
	if err != nil {
		return "", err
	}

	defer configFile.Close()

	// Decode the config file
	var config map[string]string
	err = yaml.NewDecoder(configFile).Decode(&config)
	if err != nil {
		return "", err
	}

	// Get the lock command
	lockCmd, ok := config["lock_command"]
	if !ok {
		return "", fmt.Errorf("no lock command found in config file")
	}

	supportedLockers := []string{"swaylock", "i3lock", "swaylock-blur", "swaylock-fancy"}
	if !lo.Contains(supportedLockers, lockCmd) {
		return "", fmt.Errorf("invalid lock command")
	}

	return lockCmd, nil
}
