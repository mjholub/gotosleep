package main

import (
	"os"

	"github.com/snapcore/snapd/polkit"
)

func checkAuthorization() (bool, error) {
	const actionId = "org.freedesktop.policykit.exec"

	pid := int32(os.Getpid())
	uid := uint32(os.Getuid())

	authorized, err := polkit.CheckAuthorization(pid,
		uid,
		actionId,
		nil,
		polkit.CheckAllowInteraction,
	)
	if err != nil {
		return false, err
	}

	return authorized, nil
}
