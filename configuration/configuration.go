package configuration

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

/**
Method that reads the configuration file.
If a environment variable is set, the program will read the configuration
file from the path provided otherwize it will use the path coded in hard
 */
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

/**
Method that checks that the configuration file is consistent.
If a device name that does not exist in the device list
or if a zwave device that does not exist in the zwave list
are used, error message will be displayed and the program will stop
 */
func checkConfiguration(config *Configuration) {
	exists := true;
	// Check if devices that are used in commands are in the device list
	for _, command := range config.Commands {
		for _, instruction := range command.Instructions {
			if checkDeviceExists(instruction.DeviceName, config.Devices) == false {
				fmt.Printf("ERROR : device %s does not exist in list of devices \n", instruction.DeviceName)
				exists = false
			}
		}
	}
	// Checks if the zwave name used in commands exists in the Zwave list
	// It also complets the command with the ip of the zwave device for easier use later on
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

/**
Method that checks if a zwave name exists in the list of zwave objects
 */
func checkZwaveExists(zwave string, list []Zwave) string {
	for _, value := range list {
		if zwave == value.Name {
			return value.Ip
		}
	}
	return ""
}

/**
Method that checks if a device exists in the list of devices
 */
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