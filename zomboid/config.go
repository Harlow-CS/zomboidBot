package zomboid

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	//"encoding/json"
	"io/ioutil"

	"gopkg.in/ini.v1"
)

var (

	serverConfigFilesPath = os.Getenv("server_config_files_path")
	whitelistedReadSettings = os.Getenv("whitelisted_read_settings")

)

/*
	* Gets current config, returns it as a string
	* Only returns whitelisted key values
*/
func GetCurrentConfig(serverName string) string {

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

	fileString += "===== Sandbox Settings =====\n"

	// get the Sandbox Settings as well, they're interesting and fully configurable
	// Open sandbox json file
	sandboxFilePath := filepath.Join(serverConfigFilesPath, serverName + "_sandbox.json")
	jsonFile, err := os.Open(sandboxFilePath)
	// if we os.Open returns an error then handle it
	if err != nil {
			fmt.Println(err)
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	jsonString := string(byteValue)

	fmt.Println(jsonString)
	fileString += jsonString

	return fileString

}

/*
	* Updates server-name.ini settings
*/
/*
func UpdateConfig(serverName string, newConfig map[string]interface{}) {

	// Read in current server config
	configFilePath := filepath.Join(serverConfigFilesPath, serverName + ".ini")
	currentConfig, _ := ini.Load(configFilePath)

	// set pvp
	if option, ok := newConfig["PVP"]; ok {
		currentConfig.Section("").Key("PVP").SetValue(option.StringValue())
	}

}
*/