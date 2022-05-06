package zomboid

import (
	"log"
	"os"
	"os/exec"
	//"path/filepath"
	//"encoding/json"
	//"io/ioutil"
	//"time"
	// "strings"
	"strconv"

)

// global var of tracking currently running server process
var (
	Server *os.Process = nil
)

var (

	installationPath = os.Getenv("zomboid_cli_path")

)

/*
	* Gets whether the server is active or not
*/
func IsServerActive() bool {
	output, err := exec.Command("systemctl", "is-active", "zomboid").Output()
	if err != nil {
		log.Printf("Failed to see if zomboid is active:\n%s", err)
	}

	if (string(output) == "active") {
		return true
	} else {
		return false
	}
}

/*
	* Captures the server process
*/
func GetServerProcess() {
	output, err := exec.Command("systemctl", "show", "--property", "MainPID", "zomboid").Output()
	if err != nil {
		log.Printf("Failed to get main PID of zomboid:\n%s", err)
	}

	pid, _ := strconv.Atoi(string(output))
	Server, _ = os.FindProcess(pid)

}

/*
	* Captures the server process
*/
func StartServer() {
	cmd := exec.Command("systemctl", "start", "zomboid")
	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start server:\n%s", err)
	}

}

/*
	* Captures the server process
*/
func RestartServer() {
	cmd := exec.Command("systemctl", "restart", "zomboid")
	if err := cmd.Start(); err != nil {
		log.Printf("Failed to restart server:\n%s", err)
	}

}
