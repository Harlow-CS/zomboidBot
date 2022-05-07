package zomboid

import (
	"log"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"io/ioutil"

	"gopkg.in/ini.v1"
)

var (

	serverConfigFilesPath = os.Getenv("server_config_files_path")
	whitelistedReadSettings = os.Getenv("whitelisted_read_settings")
	whitelistedWriteSettings = os.Getenv("whitelisted_write_settings")

)

/*
	* Gets current config, returns it as a string
	* Only returns whitelisted key values
*/
func GetServerConfig(serverName string) string {

	// Read in current server config
	configFilePath := filepath.Join(serverConfigFilesPath, serverName + ".ini")

	// if the file doesn't exist, return before we segfault
	if _, err := os.Stat(configFilePath); err != nil {
		return fmt.Sprintf("Server config file '%s' did not exist", configFilePath)
	}

	currentConfig, _ := ini.Load(configFilePath)

	values := currentConfig.Section("").Keys()
	keys := currentConfig.Section("").KeyStrings()

	fileString := "===== Server Settings =====\n"
	for i := 0; i < len(keys); i++ {
		key := keys[i]
		value := values[i]

		if (strings.Contains(whitelistedReadSettings, key)) {
			fileString += fmt.Sprintf("%s: %s\n", key, value)
		}
	}

	return fileString

}

/*
	* Gets current config, returns it as a string
	* Only returns whitelisted key values
*/
func GetSandboxConfig(serverName string) string {

	fileString := "===== Sandbox Settings =====\n"

	// get the Sandbox Settings as well, they're interesting and fully configurable
	// Open sandbox json file
	sandboxFilePath := filepath.Join(serverConfigFilesPath, serverName + "_sandbox.json")
	jsonFile, err := os.Open(sandboxFilePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Println(err)
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	jsonString := string(byteValue)

	fileString += jsonString

	return fileString

}

/*
	* Updates server-name.ini settings
*/
func UpdateServerConfig(serverName string, newConfig string) {

	// Read in current server config
	configFilePath := filepath.Join(serverConfigFilesPath, serverName + ".ini")
	currentConfig, _ := ini.Load(configFilePath)

	// split key/value pairs
	parameters := strings.Split(newConfig, ",")
	for i := 0; i < len(parameters); i++ {
		parameter := parameters[i]
		splitValue := strings.Split(parameter, "=")

		key := splitValue[0]
		value := splitValue[1]

		// ensure this key is in the write whitelist
		if (strings.Contains(whitelistedWriteSettings, key)) {
			currentConfig.Section("").Key(key).SetValue(value)
		}

		currentConfig.SaveTo(configFilePath)

	}

}
