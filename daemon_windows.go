// Copyright 2016 Igor Dolzhikov. All rights reserved.
// Use of this source code is governed by
// license that can be found in the LICENSE file.

// Package daemon windows version
package daemon

import (
	"fmt"
	"os/exec"
	"errors"
	"github.com/kardianos/osext"
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
	cmd := exec.Command("nssm.exe", "remove", windows.name)
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
}

// Stop the service
func (windows *windowsRecord) Stop() (string, error) {
	stopAction := "Stopping " + windows.description + ":"
	cmd := exec.Command("nssm.exe", "stop", windows.name)
	err := cmd.Run()
	if err != nil {
		return stopAction + failed, err
	}
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
	return osext.Executable()
}
