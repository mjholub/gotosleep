package main

import (
	"errors"
	"os"

	"github.com/godbus/dbus/v5"
)

func echoMemToPowerState() error {
	conn, err := dbus.SystemBus()
	if err != nil {
		return err
	}

	obj := conn.Object("org.freedesktop.PolicyKit1", "/org/freedesktop/PolicyKit1/Authority")

	subject := map[string]dbus.Variant{
		"system-bus-name": dbus.MakeVariant(conn.Names()[0]),
	}
	details := map[string]string{}

	call := obj.Call("org.freedesktop.PolicyKit1.Authority.CheckAuthorization", 0,
		[]interface{}{"system-bus-name", subject},
		"org.freedesktop.login1.hibernate",
		details,
		uint32(1), // AllowUserInteraction flag
		"",
	)

	if call.Err != nil {
		return call.Err
	}

	authorized := call.Body[0].(dbus.Variant).Value().([]interface{})[0].(bool)
	if !authorized {
		return errors.New("not authorized to suspend to memory")
	}

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
