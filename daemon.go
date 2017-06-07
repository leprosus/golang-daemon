package golang_daemon

import (
	"path/filepath"
	"os"
	"os/exec"
	"syscall"
	"fmt"
)

type Daemon struct {
	withCli bool
}

func New() Daemon {
	return Daemon{}
}

func (daemon *Daemon) StartWithCLI(mainLoop func()) (err error) {
	if len(os.Args) > 1 {
		daemon.withCli = true

		switch os.Args[1] {
		case "run":
			mainLoop()
		case "start":
			err = daemon.Start(mainLoop)
			if daemon.IsDaemonised() {
				fmt.Println("Daemon is started")
			} else {
				fmt.Fprintln(os.Stderr, "Daemon is still stopped")
			}
		case "restart":
			err = daemon.Stop()
			err = daemon.Start(mainLoop)
			if daemon.IsDaemonised() {
				println("Daemon is restarted")
			} else {
				fmt.Fprintln(os.Stderr, "Daemon is stopped. Can't run it")
			}
		case "stop":
			err = daemon.Stop()
			if daemon.IsDaemonised() {
				fmt.Fprintln(os.Stderr, "Daemon is still started")
			} else {
				fmt.Println("Daemon is stopped")
			}
		case "status":
			if daemon.IsDaemonised() {
				println("run")
			} else {
				println("stop")
			}
		default:
			help := "Usage:\n" +
				"\trun\tto run script in foreground mode\n" +
				"\tstart\tto start as daemon\n" +
				"\tstop\tto stop daemon\n" +
				"\trestart\tto restart of the daemon\n" +
				"\tstatus\treturns daemon status\n" +
				"\thelp\tto print this help\n"

			println(help)
		}
	} else {
		mainLoop()
	}

	return
}

func (daemon Daemon) Start(mainLoop func()) (err error) {
	if !daemon.IsDaemonised() {
		progName := os.Args[0]

		var path string
		if path, err = filepath.Abs(progName); err != nil {
			return
		}

		progArgs := os.Args[1:]
		if daemon.withCli &&
			progArgs[0] == "start" {
			progArgs = os.Args[2:]
		}

		cmd := exec.Command(path, progArgs...)
		if err = cmd.Start(); err != nil {
			return
		}

		os.Exit(0)
	}

	mainLoop()

	_, err = syscall.Setsid()

	return
}

func (daemon Daemon) Stop() (err error) {
	var path string
	if path, err = filepath.Abs(os.Args[0]); err == nil {
		exec.Command("pkill", "-f", path).Run()
	}

	return
}

func (daemon Daemon) IsDaemonised() bool {
	if path, err := filepath.Abs(os.Args[0]); err == nil {
		var out []byte
		out, _ = exec.Command("pgrep", "-f", path).Output()
		return len(out) > 0
	}

	return false
}
