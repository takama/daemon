// Copyright 2014 Igor Dolzhikov. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package daemon for use with Go (golang) services.
//
// Package daemon provides primitives for daemonization of golang services.
// This package is not provide implementation of user daemon,
// accordingly must have root rights to install/remove service.
// In the current implementation is only supported Linux and Mac Os X daemon.
//
// Example:
//
//	package main
//
//	import (
// 		"fmt"
//		"os"
//		"github.com/takama/daemon"
//	)
//
//	const (
//		name        = "myservice"
//		description = "Some explanation of the service purpose"
//	)
//
//	type Service struct {
//		daemon.Daemon
//	}
//
//	func (service *Service) Manage() (string, error) {
//		// if received any kind of command, do it
//		if len(os.Args) > 1 {
//			command := os.Args[1]
//			switch command {
//			case "install":
//				return service.Install()
//			case "remove":
//				return service.Remove()
//			case "start":
//				return service.Start()
//			case "stop":
//				return service.Stop()
//			case "status":
//				return service.Status()
//			}
//		}
//
//		// Do something, call your goroutines, etc
//
//		return "Usage: myservice install | remove | start | stop | status", nil
//	}
//
//	func main() {
//		srv, err := daemon.New(name, description)
//		if err != nil {
//			fmt.Println("Error: ", err)
//			os.Exit(1)
//		}
//		service := &Service{srv}
//		status, err := service.Manage()
//		if err != nil {
//			fmt.Println(status, "\nError: ", err)
//			os.Exit(1)
//		}
//		fmt.Println(status)
//	}
//
// Go daemon
package daemon

import (
	"os"
	"os/exec"
	"os/user"
)

// Service constants
const (
	rootPrivileges = "You must have root user privileges. Possibly using 'sudo' command should help"
	success        = "\t\t\t\t\t[  \033[32mOK\033[0m  ]" // Show colored "OK"
	failed         = "\t\t\t\t\t[\033[31mFAILED\033[0m]" // Show colored "FAILED"
)

// Daemon interface has standard set of a methods/commands
type Daemon interface {

	// Install the service into the system
	Install() (string, error)

	// Remove the service and all corresponded files from the system
	Remove() (string, error)

	// Start the service
	Start() (string, error)

	// Stop the service
	Stop() (string, error)

	// Status - check the service status
	Status() (string, error)
}

// New - Create a new daemon
//
// name: name ot the service, must be match with executable file name;
// description: any explanation, what is the service, its purpose
func New(name, description string) (Daemon, error) {
	return newDaemon(name, description)
}

// Lookup path for executable file
func executablePath(name string) (string, error) {
	if path, err := exec.LookPath(name); err == nil {
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			return execPath()
		}
		return path, nil
	}
	return execPath()
}

// Check root rights to use system service
func checkPrivileges() bool {

	if user, err := user.Current(); err == nil && user.Gid == "0" {
		return true
	}
	return false
}
