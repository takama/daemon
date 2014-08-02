// Copyright 2014 Igor Dolzhikov. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package daemon windows version
package daemon

import (
	"errors"
)

// WindowsRecord - standard record (struct) for windows version of daemon package
type WindowsRecord struct {
	name        string
	description string
}

func newDaemon(name, description string) (*WindowsRecord, error) {

	return &WindowsRecord{name, description}, nil
}

// Install the service
func (windows *WindowsRecord) Install() (string, error) {
	installAction := "Install " + windows.description + ":"

	return installAction + failed, errors.New("windows daemon not supported")
}

// Remove the service
func (windows *WindowsRecord) Remove() (string, error) {
	removeAction := "Removing " + windows.description + ":"

	return removeAction + failed, errors.New("windows daemon not supported")
}

// Start the service
func (windows *WindowsRecord) Start() (string, error) {
	startAction := "Starting " + windows.description + ":"

	return startAction + failed, errors.New("windows daemon not supported")
}

// Stop the service
func (windows *WindowsRecord) Stop() (string, error) {
	stopAction := "Stopping " + windows.description + ":"

	return stopAction + failed, errors.New("windows daemon not supported")
}

// Status - Get service status
func (windows *WindowsRecord) Status() (string, error) {

	return "Status could not defined", errors.New("windows daemon not supported")
}
