package zomboid

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"encoding/json"
	"io/ioutil"
	"time"
	"strings"

	"gopkg.in/ini.v1"
)

// global var of tracking currently running server process
var (
	Server *os.Process = nil
)

var (

	installationPath = os.Getenv("zomboid_cli_path")
	serverConfigFilesPath = os.Getenv("server_config_files_path")

)

/*
	* Gets whether the server is active or not
*/
func IsServerActive() bool {
	cmd := exec.Command("systemctl", "is-active", "zomboid")
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to start server:\n%s", err)
	}

	output := strings.TrimSpace(string(cmd))
	if (output == "active") {
		return true
	} else {
		return false
	}
}

/*
	* Captures the server process
*/
func GetServerProcess() {
	cmd := exec.Command("systemctl", "show", "--property", "MainPID", "zomboid")
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to start server:\n%s", err)
	}

	pid := strings.TrimSpace(string(cmd))
	Server = os.FindProcess(pid)

}
