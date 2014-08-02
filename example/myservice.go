package main

import (
	"fmt"
	"os"

	"github.com/takama/daemon"
)

const (
	name        = "myservice"
	description = "Some explanation of the service purpose"
)

type Service struct {
	daemon.Daemon
}

func (service *Service) Manage() (string, error) {
	// if received any kind of command, do it
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		}
	}

	// Do something, call your goroutines, etc

	return "Usage: myservice install | remove | start | stop | status", nil
}

func main() {
	srv, err := daemon.New(name, description)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	service := &Service{srv}
	status, err := service.Manage()
	if err != nil {
		fmt.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)

}
