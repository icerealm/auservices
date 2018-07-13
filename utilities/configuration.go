package utilities

import (
	"encoding/json"
	"os"
)

//Configuration represents application configuration
type Configuration struct {
	ApplicationPort   int
	MsgURL            string
	MsgClusterID      string
	DbURL             string
	CategoryChannelID string
}

//appConfig represent configuration
var appConfig *Configuration

//GetConfiguration return configuration copy.
func GetConfiguration() Configuration {
	return *appConfig
}

//LoadConfiguration load configuration struct
func LoadConfiguration(filename string) (*Configuration, error) {
	appConfig = &Configuration{}
	env := os.Getenv("AppEnvironment")
	err := loadFile(filename, appConfig)
	if env == "PROD" {
		initProdConfiguration(appConfig)
		return appConfig, err
	}
	return appConfig, err
}

func loadFile(filename string, cfg *Configuration) error {
	fd, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fd.Close()

	decoder := json.NewDecoder(fd)
	return decoder.Decode(&cfg)
}

func initProdConfiguration(config *Configuration) {
	config.MsgURL = os.Getenv("MsgURL")
	config.DbURL = os.Getenv("DbURL")
}
