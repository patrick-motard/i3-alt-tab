package main

import (
	"fmt"
	"github.com/takama/daemon"
	"go.i3wm.org/i3"
	"log"
	"net"
	"os"
)

func server(c net.Conn) {
	for {
		buf := make([]byte, 512)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}
		data := buf[0:nr]
		println("Server got:", string(data))
		_, err = c.Write(data)
		if err != nil {
			log.Fatal("Write: ", err)
		}
	}
}

type Service struct {
	daemon.Daemon
}

func (service *Service) Manage() (string, error) {
	usage := "Usage: myservice install | remove | start | stop | status"
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			service.Start()
			l, err := net.Listen("unix", "/tmp/truck.sock")
			if err != nil {
				log.Fatal("listen error:", err)
			}
			for {
				fd, err := l.Accept()
				if err != nil {
					log.Fatal("accept error:", err)
				}
				go server(fd)

			}
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}
}

func main() {

	service, err := daemon.New("truck", "description")
	if err != nil {
		log.Fatal("Error: ", err)
	}

	service.Remove()
	status, err := service.Install()

	if err != nil {
		log.Fatal(status, "\nError: ", err)
	}

	// k,err:= service.Status()
	fmt.Println(service.Status())
	fmt.Println(status)
	service.Start()

	fmt.Println(service.Status())
	fmt.Println("hello world")
	// tree, err := i3.GetTree()
	tree, err := i3.GetWorkspaces()

	if err != nil {
		log.Fatal(err)
	}
	for i, v := range tree {
		fmt.Println(i, v)
	}
	i3.RunCommand("workspace ")
	// fmt.Printf("%+v\n", tree[0])
	// fmt.Println(tree[0])
}
