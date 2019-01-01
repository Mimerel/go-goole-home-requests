package configuration

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func ReadConfiguration() (*Configuration) {
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

	var config *Configuration

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	} else {
		checkConfiguration(config)
		fmt.Printf("Configuration Loaded : %+v \n", config)
	}
	return config
}

func checkConfiguration(config *Configuration) {
	exists := true;
	for _, command := range config.Commands {
		for _, instruction := range command.Instructions {
			if checkDeviceExists(instruction.DeviceName, config.Devices) == false {
				fmt.Printf("ERROR : device %s does not exist in list of devices \n", instruction.DeviceName)
				exists = false
			}
		}
	}
	for i:=0; i< len(config.Devices); i++ {
		ip := checkZwaveExists(config.Devices[i].Zwave, config.Zwaves)
		if ip == "" {
			fmt.Printf("ERROR : zwave %s does not exist in list of zwave devices \n", config.Devices[i].Zwave)
			exists = false
		} else {
			config.Devices[i].Url = ip
		}
	}
	if exists == false {
		os.Exit(1)
	}
}

func checkZwaveExists(zwave string, list []Zwave) string {
	for _, value := range list {
		if zwave == value.Name {
			return value.Ip
		}
	}
	return ""
}

func checkDeviceExists(device string, list []Device) bool {
	exists := false
	for _, value := range list {
		if device == value.Name {
			exists = true
			break
		}
	}
	return exists
}