package configurations

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

/*
	Look at moving to viper in the future
	github.com/spf13/viper
*/

var configFileName = func() string {
	if value, ok := os.LookupEnv("CONFIG_FILE"); ok {
		return value
	}
	return "config.json"
}()

type Configuration struct {
	GameServer Server `json:"gameServer"`
}

type Server struct {
	Host string `json:"host"`
	Path string `json:"path"`
}

var Config = parseConfig(configFileName)

func parseConfig(fileName string) *Configuration {
	configFile, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Failed to open config file '%v'", fileName)
	}
	defer configFile.Close()
	configData, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Fatalf("Error while reading config file '%s'", fileName)
	}

	config := Configuration{}
	if err = json.Unmarshal(configData, &config); err != nil {
		log.Fatalf("Error unmarshalling config file '%s'", fileName)
	}
	return &config
}
