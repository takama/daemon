package main

import (
	"fmt"
	"github.com/robfig/cron"
	"github.com/takama/daemon"
	"log"
	"os"
)

const (
	// name of the service
	name        = "goCron"
	description = "Cron service example"
)

var stdlog, errlog *log.Logger

// Service is the daemon service struct
type Service struct {
	d daemon.Daemon
}

func startCron(c *cron.Cron) {
	// Run 1x every min
	c.AddFunc("* * * * * *", func() { makeFile() })
	c.Start()
	for {

	}
}

var times int

func makeFile() {
	// create a simple file $times.txt
	times++
	f, err := os.Create(fmt.Sprintf("/SOMEPATH/%d.txt", times))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
}

// Manage by daemon commands or run the daemon
func (service *Service) Manage() (string, error) {

	// Create a new cron manager
	c := cron.New()
	usage := "Usage: cronStock install | remove | start | stop | status"
	// If received any kind of command, do it
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.d.Install()
		case "remove":
			return service.d.Remove()
		case "start":
			return service.d.Start()
		case "stop":
			// No need to explicitly stop cron since job will be killed
			return service.d.Stop()
		case "status":
			return service.d.Status()
		default:
			return usage, nil
		}
	}
	// Begin cron job
	go startCron(c)
	for {
		//EventLoop to keep cron running
	}
	// Unreachable, but required
	return usage, nil
}
func init() {
	stdlog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	errlog = log.New(os.Stderr, "", log.Ldate|log.Ltime)
}
func main() {
	srv, err := daemon.New(name, description)
	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}
	service := &Service{srv}
	status, err := service.Manage()
	if err != nil {
		errlog.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}
