package zomboid

import (
	"path/filepath"
	"log"
	"os"
	"os/exec"
)

// global var of tracking currently running server process
var (
	Server *os.Process = nil
	installationPath = os.Getenv("zomboid_cli_path")
)

/*
	* Gets whether the server is active or not
*/
func IsServerActive() bool {
	return Server != nil
}

/*
	* Captures the server process
*/
func StartServer() {
	serverExecPath := filepath.Join(installationPath, "start-server.sh")
	cmd := exec.Command(serverExecPath)
	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start server:\n%s", err)
	}
	Server = cmd.Process
}

/*
	* Kills the server
*/
func StopServer() {
	Server.Kill()
	Server = nil
}
