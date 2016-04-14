// Copyright 2016 Igor Dolzhikov. All rights reserved.
// Use of this source code is governed by
// license that can be found in the LICENSE file.

// Package daemon windows version
package daemon

import (
	"errors"
)

var ErrWindowsUnsupported = errors.New("windows daemon is not supported")

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

	return installAction + failed, ErrWindowsUnsupported
}

// Remove the service
func (windows *windowsRecord) Remove() (string, error) {
	removeAction := "Removing " + windows.description + ":"

	return removeAction + failed, ErrWindowsUnsupported
}

// Start the service
func (windows *windowsRecord) Start() (string, error) {
	startAction := "Starting " + windows.description + ":"

	return startAction + failed, ErrWindowsUnsupported
}

// Stop the service
func (windows *windowsRecord) Stop() (string, error) {
	stopAction := "Stopping " + windows.description + ":"

	return stopAction + failed, ErrWindowsUnsupported
}

// Status - Get service status
func (windows *windowsRecord) Status() (string, error) {

	return "Status could not defined", ErrWindowsUnsupported
}

// Get executable path
func execPath() (string, error) {
	return "", ErrWindowsUnsupported
}
