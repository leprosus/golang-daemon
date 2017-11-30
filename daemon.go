package golang_daemon

import (
	"path/filepath"
	"os"
	"os/exec"
	"syscall"
	"strings"
	"time"
	"io/ioutil"
)

type Script struct {
	abs  string
	name string
}

type Daemon struct {
	withCli    bool
	allowUsers []string
	script     Script
	mainLoop   func()
	pidPath    string
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

func New(mainLoop func()) Daemon {
	return Daemon{
		mainLoop: mainLoop,
		script:   Script{}}
}

func (daemon *Daemon) AllowUser(userName string) {
	daemon.allowUsers = append(daemon.allowUsers, userName)
}

func (daemon *Daemon) PIDFile(pidPath string) {
	daemon.pidPath = pidPath
}

func (daemon *Daemon) StartWithCLI() (err error) {
	if !daemon.IsAllowExec() {
		println(havntAccess)

		return
	}

	if len(os.Args) > 1 {
		err = daemon.splitScriptName()
		if err != nil {
			return
		}

		if os.Args[1] == "status" {
			if daemon.IsDaemonised() {
				println(daemon.script.name, "is running")
			} else {
				println(daemon.script.name, "is stopped")
			}

			return
		}

		daemon.withCli = true

		switch os.Args[1] {
		case "run":
			daemon.mainLoop()
		case "start":
			err = daemon.doStart()
			if err == nil {
				println(daemon.script.name, "is running")
			}
		case "restart":
			err = daemon.doStop()
			if err == nil {
				println(daemon.script.name, "is stopped")
			}

			for daemon.IsDaemonised() {
				time.Sleep(time.Second)
			}

			err = daemon.doStart()
			if err == nil {
				println(daemon.script.name, "is running")
			}
		case "stop":
			err = daemon.doStop()
			if err == nil {
				println(daemon.script.name, "is stopped")
			}
		default:
			println(help)
		}
	} else {
		daemon.mainLoop()
	}

	return
}

func (daemon *Daemon) splitScriptName() (err error) {
	daemon.script.abs, err = filepath.Abs(os.Args[0])
	if err != nil {
		return
	}

	info, err := os.Stat(daemon.script.abs)
	if err != nil {
		return
	}

	daemon.script.name = info.Name()

	return
}

func (daemon Daemon) doStart() (err error) {
	if !daemon.IsDaemonised() {
		err = daemon.Start()
	}

	return
}

func (daemon Daemon) doStop() (err error) {
	if daemon.IsDaemonised() {
		err = daemon.Stop()
	}

	return
}

func (daemon Daemon) Start() (err error) {
	if !daemon.IsDaemonised() {
		progArgs := os.Args[1:]
		if daemon.withCli &&
			(progArgs[0] == "start" || progArgs[0] == "restart") {
			progArgs = os.Args[2:]
		}

		cmd := exec.Command(daemon.script.abs, progArgs...)
		if err = cmd.Start(); err != nil {
			return
		}

		if len(daemon.pidPath) > 0 {
			var out []byte
			out, _ = exec.Command("pgrep", "-f", daemon.script.abs).Output()
			ioutil.WriteFile(daemon.pidPath, out, 0440)
		}

		return
	}

	_, err = syscall.Setsid()

	return
}

func (daemon Daemon) Stop() (err error) {
	exec.Command("pkill", "-f", daemon.script.abs).Output()

	if len(daemon.pidPath) > 0 {
		os.Remove(daemon.pidPath)
	}

	return
}

func (daemon Daemon) IsDaemonised() bool {
	var out []byte
	out, _ = exec.Command("pgrep", "-f", daemon.script.abs).Output()

	return len(out) > 0
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
