package daemon

import (
	"os"
	"os/exec"
	"log"
"strings"
)

// Service constants
const (
	rootPrivileges = "You must have root user privileges. Possibly using 'sudo' command should help"
	success        = "\t\t\t\t\t[  \033[32mOK\033[0m  ]" // Show colored "OK"
	failed         = "\t\t\t\t\t[\033[31mFAILED\033[0m]" // Show colored "FAILED"
)

// Lookup path for executable file
func executablePath(name string) (string, error) {
	if path, err := exec.LookPath(name); err == nil {
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			return execPath()
		}
		return path, nil
	}
	return execPath()
}

// Check root rights to use system service
func checkPrivileges() bool {
	cmd := exec.Command("id", "-g")
	if output, err := cmd.Output(); err == nil {
		gid := strings.TrimSpace(string(output))
		return gid == "0"
	} else {
		log.Println(err)
	}
	return false
}
