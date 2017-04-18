// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by
// license that can be found in the LICENSE file.

// Package daemon windows version
package daemon

import (
	"errors"
	"fmt"
	"os/exec"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

var ErrWindowsUnsupported = errors.New("Adding as a service failed. Download and place nssm.exe in the path to install this service as an service. NSSM url: https://nssm.cc/")

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
	adminAccessNecessary := "Administrator access is needed to install a service."

	execp, err := execPath()

	if err != nil {
		return installAction + failed, err
	}

	cmd := exec.Command("nssm.exe", "install", windows.name, execp)
	out, err := cmd.Output()
	if err != nil {
		if len(out) > 0 {
			fmt.Println(string(out))
		} else {
			fmt.Println("No output. Probably service already exists. Try uninstall first.")
		}
		return installAction + failed, err
	}
	if len(out) == 0 {
		return adminAccessNecessary, errors.New(adminAccessNecessary)
	}
	return installAction + " completed.", nil
}

// Remove the service
func (windows *windowsRecord) Remove() (string, error) {
	removeAction := "Removing " + windows.description + ":"
	cmd := exec.Command("nssm.exe", "remove", windows.name, "confirm")
	err := cmd.Run()
	if err != nil {
		return removeAction + failed, err
	}
	return removeAction + " completed.", nil
}

// Start the service
func (windows *windowsRecord) Start() (string, error) {
	startAction := "Starting " + windows.description + ":"
	cmd := exec.Command("nssm.exe", "start", windows.name)
	err := cmd.Run()
	if err != nil {
		return startAction + failed, err
	}
	return startAction + " completed.", nil
}

// Stop the service
func (windows *windowsRecord) Stop() (string, error) {
	stopAction := "Stopping " + windows.description + ":"
	cmd := exec.Command("nssm.exe", "stop", windows.name)
	err := cmd.Run()
	if err != nil {
		return stopAction + failed, err
	}
	return stopAction + " completed.", nil
}

// Status - Get service status
func (windows *windowsRecord) Status() (string, error) {
	cmd := exec.Command("nssm.exe", "status", windows.name)
	out, err := cmd.Output()
	if err != nil {
		return "Getting status:" + failed, err
	}
	return "Status: " + string(out), nil
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
