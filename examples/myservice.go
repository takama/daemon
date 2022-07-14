// Example of a daemon with echo service
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/takama/daemon"
)

const (

	// name of the service
	name        = "myservice"
	description = "My Echo Service"

	// port which daemon should be listen
	port = ":9977"
)

// dependencies that are NOT required by the service, but might be used
var dependencies = []string{ /*"dummy.service"*/ }

var stdlog, errlog *log.Logger

// MyService implements the daemon.Executable interface
// and represents the actual service behavior
type MyService struct {
	listen chan net.Conn
}

// Start gets the service up
func (mysvc *MyService) Start() {
	// Set up listener for defined host and port
	listener, err := net.Listen("tcp", port)
	if err != nil {
		errlog.Println("Possibly was a problem with the port binding", err)
		return
	}

	// set up channel on which to send accepted connections
	mysvc.listen = make(chan net.Conn, 100)
	go acceptConnection(listener, mysvc.listen)

	// loop work cycle with accept connections or interrupt
	// by system signal
	go func() {
		for {
			select {
			case conn, ok := <-mysvc.listen:
				if !ok {
					stdlog.Println("Closing connections")
					listener.Close()
					return
				}
				go handleClient(conn)
			}
		}
	}()
}

// Stop shuts down the service
func (mysvc *MyService) Stop() {
	close(mysvc.listen)
}

// Run is invoked when the service is run in interective mode
// (ie during development). On Windows it is never invoked
func (mysvc *MyService) Run() {
	mysvc.Start()
	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	// loop work cycle with accept connections or interrupt
	// by system signal
loop:
	for {
		select {
		case killSignal := <-interrupt:
			stdlog.Println("Got signal:", killSignal)
			if killSignal == os.Interrupt {
				stdlog.Println("Daemon was interrupted by system signal")
			}
			stdlog.Println("Daemon was killed")
			break loop
		}
	}

	mysvc.Stop()
}

// Service has embedded daemon
// daemon.Daemon abstracts os specific service mechanics
type Service struct {
	daemon.Daemon
}

// Manage by daemon commands or run the daemon
func (service *Service) Manage() (string, error) {

	usage := "Usage: myservice install | remove | start | stop | status"

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
		default:
			return usage, nil
		}
	}

	mysvc := &MyService{}
	return service.Run(mysvc)
}

// Accept a client connection and collect it in a channel
func acceptConnection(listener net.Listener, listen chan<- net.Conn) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		listen <- conn
	}
}

func handleClient(client net.Conn) {
	for {
		buf := make([]byte, 4096)
		numbytes, err := client.Read(buf)
		if numbytes == 0 || err != nil {
			return
		}
		client.Write(buf[:numbytes])
	}
}

func init() {
	stdlog = log.New(os.Stdout, "", 0)
	errlog = log.New(os.Stderr, "", 0)
}

func main() {
	daemonKind := daemon.SystemDaemon
	if runtime.GOOS == "darwin" {
		daemonKind = daemon.UserAgent
	}
	srv, err := daemon.New(name, description, daemonKind, dependencies...)
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
