package golang_daemon

import (
	"path/filepath"
	"os"
	"os/exec"
	"syscall"
	"strings"
)

type Daemon struct {
	withCli    bool
	allowUsers []string
}

const (
	help = "Usage:\n" +
		"\trun\tto run script in foreground mode\n" +
		"\tstart\tto start as daemon\n" +
		"\tstop\tto stop daemon\n" +
		"\trestart\tto restart of the daemon\n" +
		"\tstatus\treturns daemon status\n" +
		"\thelp\tto print this help\n"

	havntAccess = "You havn't access to execute"
)

func New() Daemon {
	return Daemon{}
}

func (daemon *Daemon) AllowUser(userName string) {
	daemon.allowUsers = append(daemon.allowUsers, userName)
}

func (daemon *Daemon) StartWithCLI(mainLoop func()) (err error) {
	if len(os.Args) > 1 {
		daemon.withCli = true

		switch os.Args[1] {
		case "run":
			if daemon.IsAllowExec() {
				mainLoop()
			} else {
				println(havntAccess)
			}
		case "start":
			if daemon.IsAllowExec() {
				err = daemon.Start(mainLoop)
			} else {
				println(havntAccess)
			}
		case "restart":
			if daemon.IsAllowExec() {
				daemon.Stop()
				err = daemon.Start(mainLoop)
			} else {
				println(havntAccess)
			}
		case "stop":
			if daemon.IsAllowExec() {
				err = daemon.Stop()
			} else {
				println(havntAccess)
			}
		case "status":
			if daemon.IsDaemonised() {
				println("run")
			} else {
				println("stop")
			}
		default:
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
		exec.Command("pkill", "-f", path).Output()
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

func (daemon Daemon) IsAllowExec() bool {
	if len(daemon.allowUsers) == 0 {
		return true
	}

	bytes, err := exec.Command("whoami", ).Output()
	if err != nil {
		return false
	}

	curUser := strings.TrimSpace(string(bytes))

	for _, user := range daemon.allowUsers {
		if curUser == user {
			return true
		}
	}

	return false
}
