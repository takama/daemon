// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by
// license that can be found in the LICENSE file.

// Package daemon windows version
package daemon

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

// windowsRecord - standard record (struct) for windows version of daemon package
type windowsRecord struct {
	name         string
	description  string
	dependencies []string
}

func newDaemon(name, description string, dependencies []string) (Daemon, error) {

	return &windowsRecord{name, description, dependencies}, nil
}

// Install the service
func (windows *windowsRecord) Install(args ...string) (string, error) {
	installAction := "Install " + windows.description + ":"

	execp, err := execPath()

	if err != nil {
		return installAction + failed, err
	}

	cmdArgs := []string{"create", windows.name, "start=auto", "binPath=" + execp}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("sc", cmdArgs...)
	_, err = cmd.Output()
	if err != nil {
		return installAction + failed, getWindowsError(err)
	}
	return installAction + " completed.", nil
}

// Remove the service
func (windows *windowsRecord) Remove() (string, error) {
	removeAction := "Removing " + windows.description + ":"
	cmd := exec.Command("sc", "delete", windows.name, "confirm")
	err := cmd.Run()
	if err != nil {
		return removeAction + failed, getWindowsError(err)
	}
	return removeAction + " completed.", nil
}

// Start the service
func (windows *windowsRecord) Start() (string, error) {
	startAction := "Starting " + windows.description + ":"
	cmd := exec.Command("sc", "start", windows.name)
	err := cmd.Run()
	if err != nil {
		return startAction + failed, getWindowsError(err)
	}
	return startAction + " completed.", nil
}

// Stop the service
func (windows *windowsRecord) Stop() (string, error) {
	stopAction := "Stopping " + windows.description + ":"
	cmd := exec.Command("sc", "stop", windows.name)
	err := cmd.Run()
	if err != nil {
		return stopAction + failed, err
	}
	return stopAction + " completed.", nil
}

// Status - Get service status
func (windows *windowsRecord) Status() (string, error) {
	cmd := exec.Command("sc", "query", windows.name)
	out, err := cmd.Output()
	if err != nil {
		return "Getting status:" + failed, getWindowsError(err)
	}
	return "Status: " + "SERVICE_" + getWindowsServiceState(out), nil
}

// Get executable path
func execPath() (string, error) {
	var n uint32
	b := make([]uint16, syscall.MAX_PATH)
	size := uint32(len(b))

	r0, _, e1 := syscall.MustLoadDLL(
		"kernel32.dll",
	).MustFindProc(
		"GetModuleFileNameW",
	).Call(0, uintptr(unsafe.Pointer(&b[0])), uintptr(size))
	n = uint32(r0)
	if n == 0 {
		return "", e1
	}
	return string(utf16.Decode(b[0:n])), nil
}

// Get windows error
func getWindowsError(inputError error) error {
	if exiterr, ok := inputError.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			if sysErr, ok := WinErrCode[status.ExitStatus()]; ok {
				return errors.New(fmt.Sprintf("\n %s: %s \n %s", sysErr.Title, sysErr.Description, sysErr.Action))
			}
		}
	}

	return inputError
}

// Get windows service state
func getWindowsServiceState(out []byte) string {
	regex := regexp.MustCompile("STATE.*: (?P<state_code>[0-9])  (?P<state>.*) ")
	service := regex.FindAllStringSubmatch(string(out), -1)[0]

	return service[2]
}
