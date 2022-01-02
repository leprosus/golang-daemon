package daemon

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

var cli, pid string

func Init(name string, params map[string]interface{}, pidPath string) (err error) {
	cli, err = filepath.Abs(name)
	if err != nil {
		return
	}

	pid, err = filepath.Abs(pidPath)
	if err != nil {
		return
	}

	err = os.MkdirAll(filepath.Dir(pid), 0755)
	if err != nil {
		return
	}

	if len(params) == 0 {
		return
	}

	var lines []string
	for key, val := range params {
		lines = append(lines, fmt.Sprintf("--%s=%v", key, val))
	}

	sort.Strings(lines)

	cli += " " + strings.Join(lines, " ")

	return
}

func IsRun() (ok bool) {
	out, _ := exec.Command("sh", "-c", fmt.Sprintf("pgrep -f '%s'", cli)).Output()

	return len(out) > 0
}

func Start() (err error) {
	if IsRun() {
		return
	}

	cmd := exec.Command("sh", "-c", cli)
	err = cmd.Start()
	if err != nil {
		return
	}

	return CreatePIDFile()
}

func Stop() (err error) {
	_ = exec.Command("sh", "-c", fmt.Sprintf("pkill -f '%s'", cli)).Run()

	return RemovePIDFile()
}

func CreatePIDFile() (err error) {
	return ioutil.WriteFile(pid, []byte(strconv.Itoa(os.Getpid())), 0640)
}

func RemovePIDFile() (err error) {
	return os.Remove(pid)
}
