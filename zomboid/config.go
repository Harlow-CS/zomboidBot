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

var (

	serverConfigFilesPath = os.Getenv("server_config_files_path")

)

/*
	* Updates server-name.ini settings
*/
func UpdateServerConfig(serverName string, newConfig interface{}) {

	// Read in current server config
	configFilePath = filepath.Join(serverConfigFilesPath, serverName + ".ini")
	currentConfig, err := ini.Load(configFilePath)

	// set pvp

	// set PlayerBumpPlayer

	// StarterKit


}