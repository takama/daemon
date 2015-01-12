package daemon

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

// Service Daemon
type ServiceDaemon struct {
	Daemon
}

func (daemon *ServiceDaemon) Manage(service Service) (string, error) {

	usage := "Usage: myservice install | remove | start | stop | status"

	// if received any kind of command, do it
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return daemon.Install()
		case "remove":
			return daemon.Remove()
		case "start":
			return daemon.Start()
		case "stop":
			return daemon.Stop()
		case "status":
			return daemon.Status()
		default:
			return usage, nil
		}
	}

	process := service.GetProcess()
	process()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Set up listener for defined host and port
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(service.GetPort()))
	if err != nil {
		return "Possibly was a problem with the port binding", err
	}

	// set up channel on which to send accepted connections
	listen := make(chan net.Conn, 100)
	go acceptConnection(listener, listen)

	// loop work cycle with accept connections or interrupt
	// by system signal
	for {
		select {
		case conn := <-listen:
			go handleClient(conn)
		case killSignal := <-interrupt:
			log.Println("Got signal:", killSignal)
			log.Println("Stoping listening on ", listener.Addr())
			listener.Close()
			if killSignal == os.Interrupt {
				return "Daemon was interruped by system signal", nil
			}
			return "Daemon was killed", nil
		}
	}

	// never happen, but need to complete code
	return usage, nil
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
		client.Write(buf)
	}
}

// Services

type Service struct {
	Name        string
	Description string
	Port        int
	Process     func()
}

func (s Service) GetName() string {
	return s.Name
}

func (s Service) GetPort() int {
	return s.Port
}

func (s Service) GetDescription() string {
	return s.Description
}

func (s Service) GetProcess() func() {
	return s.Process
}

func (service Service) Daemon() {

	srv, err := New(service.GetName(), service.GetDescription())

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	service_daemon := &ServiceDaemon{srv}
	status, err := service_daemon.Manage(service)

	if err != nil {
		fmt.Println(status, "\nError: ", err)
		os.Exit(1)
	}

	fmt.Println(status)
}
