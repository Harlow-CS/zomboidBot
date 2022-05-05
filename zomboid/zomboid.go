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
)

// global var of tracking currently running server process
var (
	Server *os.Process = nil
)

var (

	installationPath = os.Getenv("zomboid_cli_path")

)

func PingServer(saveFile string) {
	log.Println("Starting server...")

	var saveFilePath = filepath.Join(savesPath, saveFile + ".zip")

	cmd := exec.Command(serverPath, "--start-server", saveFilePath, "--server-settings", settingsPath, "--bind", "0.0.0.0")
	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start server:\n%s", err)
	}

	// capture the server process
	Server = cmd.Process

	// poll and wait for the server log to say it's ready for hosting
	// if it takes > 5 minutes, exit and say it's bad
	var ready = false
	var totalSleep = 0
	for (!ready || totalSleep >= 300) {
		b, err := ioutil.ReadFile(serverCurrentLogPath) 
		if err != nil {
			log.Print(err)
		}
		serverLogs := string(b)

		if (strings.Contains(serverLogs, "Matching server game")) {
			ready = true
		} else {
			time.Sleep(10 * time.Second)
			totalSleep = totalSleep + 10
		}

	}

	if(ready) {
		log.Println("Server started")
	} else {
		log.Println("Server failed to start")
		Server.Kill()
		Server = nil
	}

}
