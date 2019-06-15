package configuration

import (
	"fmt"
	"go-goole-home-requests/models"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

/**
Method that reads the configuration file.
If a environment variable is set, the program will read the configuration
file from the path provided otherwize it will use the path coded in hard
 */
func ReadConfiguration() (*models.Configuration) {
	pathToFile := os.Getenv("LOGGER_CONFIGURATION_FILE")
	if _, err := os.Stat("/home/pi/go/src/go-goole-home-requests/configuration.yaml"); !os.IsNotExist(err) {
		pathToFile = "/home/pi/go/src/go-goole-home-requests/configuration.yaml"
	} else {
		pathToFile = "./configuration.yaml"
	}
	yamlFile, err := ioutil.ReadFile(pathToFile)

	if err != nil {
		panic(err)
	}

	var config *models.Configuration

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	} else {
		executeGoogleAssistantConfiguration(config)
		fmt.Printf("Configuration Loaded : %+v \n", config)
	}

	return config
}
