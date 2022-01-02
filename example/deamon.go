//go:build !windows
// +build !windows

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	daemon "github.com/leprosus/golang-daemon"
)

func main() {
	err := daemon.Init(os.Args[0], map[string]interface{}{}, "./app.pid")
	if err != nil {
		return
	}

	var cmd string
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "start":
		err = daemon.Start()
	case "stop":
		err = daemon.Stop()
	case "restart":
		err = daemon.Stop()
		err = daemon.Start()
	case "status":
		status := "stopped"
		if daemon.IsRun() {
			status = "started"
		}

		fmt.Printf("Application is %s\n", status)

		return
	case "":
		fallthrough
	default:
		mainLoop()
	}
}

func mainLoop() {
	for {
		log.Println("I'm daemon")
		time.Sleep(time.Minute)
	}
}
