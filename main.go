package main

import (
	"fmt"
	"github.com/op/go-logging"
	"github.com/Mimerel/go-logger-client"
	"go-goole-home-requests/configuration"
	"go-goole-home-requests/google_talk"
	"go-goole-home-requests/utils"
	"net/http"
	"os"
	"strings"
	"time"
)

var log = logging.MustGetLogger("default")

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{color:reset} %{message}`,
)

/**
Method that sends a request to a domotic zwave server to run an instruction
 */
func ExecuteRequest(url string, id string, instance string, commandClass string, level string) (err error) {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	postingUrl := "http://" + url + ":8083/ZWaveAPI/Run/devices[" + id + "].instances[" + instance + "].commandClasses[" + commandClass + "].Set(" + level + ")"
	log.Info("Request posted : %s", postingUrl)

	_, err = client.Get(postingUrl)
	if err != nil {
		fmt.Printf("Failed to execute request %s \n", postingUrl, err)
		return err
	}
	return nil
}

/**
Main method
 */
func main() {
	config := configuration.ReadConfiguration()

	// Speak text on Google Home.

	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backend)
	backendLeveled.SetLevel(logging.NOTICE, "")
	logging.SetBackend(backendLeveled, backendFormatter)
	log.Info("Application Starting")

	http.HandleFunc("/question/", func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path
		urlParams := strings.Split(urlPath, "/")
		log.Info("Request received question %s / %d", urlPath, len(urlParams))

		if len(urlParams) == 4 {
			log.Info("Request succeeded")
			AnalyseQuestionRequest(w, r, urlParams, config)
		} else {
			log.Info("Request failed")
			w.WriteHeader(500)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path
		urlParams := strings.Split(urlPath, "/")
		log.Info("Request received question %s / %d", urlPath, len(urlParams))

		if len(urlParams) == 3 {
			log.Info("Request succeeded")
			AnalyseRequest(w, r, urlParams, config)
		} else {
			log.Info("Request failed")
			w.WriteHeader(500)
		}
	})

	err := http.ListenAndServe(":9998", nil)
	if err != nil {
		log.Errorf("error %+v", err)
	}
}

/**
Method that searches for the ip(s) concerned by a room.
When an instruction is used, it will always be linked to a room
 */
func findIpOfGoogleHome(googleList []configuration.GoogleDetails, concernedRoom string) ([]string) {
	ips := []string{}
	for _, google := range googleList {
		if google.Name == concernedRoom {
			ips = google.Ip
		}
	}
	return ips
}

/**
Method that splits the instruction into an action and a instruction
 */
func getActionAndInstruction(instruction string) (action string, newInstruction string) {
	instruction = utils.ConvertInstruction(instruction)
	log.Info("instructions: <%s> ", instruction)
	mainAction := strings.Split(instruction, " ")[0]
	instruction = strings.Replace(instruction, mainAction, "", 1)
	instruction = strings.Trim(instruction, " ")
	return mainAction, instruction
}

/**
Method that checks if the action demanded exists and retrieves the information linked to this action.
 */
func checkActionValidity(actions []configuration.ActionDetails, mainAction string) (bool, string, string, string) {
	found := false
	level := ""
	actionType := ""
	for _, action := range actions {
		for _, actionName := range action.Name {
			if actionName == mainAction {
				level = action.Value
				mainAction = action.Replacement
				actionType = action.Type
				found = true
			}
		}
	}
	return found, mainAction, level, actionType
}

/**
Method that searches throw the list of possible commands for the
command sent by google home.
It first tries to find the corresponding "sentence" in its database.
IF it is found, it will check if the action is autorized in that room
If so, it will execute the command
 */
func RunDomoticCommand(config *configuration.Configuration, instruction string, concernedRoom string, mainAction string, level string) (bool) {
	found := false
	for _, ListInstructions := range config.Commands {
		for _, word := range ListInstructions.Words {
			if utils.CompareWords(word, instruction) &&
				utils.CompareRooms(ListInstructions.Rooms, concernedRoom) &&
				utils.CompareActions(ListInstructions.Actions, mainAction) {
				for _, instruction := range ListInstructions.Instructions {
					if instruction.Value == "" {
						ExecuteAction(config, level, instruction.DeviceName)

					} else {
						ExecuteAction(config, instruction.Value, instruction.DeviceName)
					}
				}
				found = true;
				break
			}
		}
	}
	return found
}

func AnalyseRequest(w http.ResponseWriter, r *http.Request, urlParams []string, config *configuration.Configuration) {
	concernedRoom := urlParams[1]
	ips := findIpOfGoogleHome(config.Googles, concernedRoom)
	if len(ips) == 0 {
		w.WriteHeader(500)
	} else {
		mainAction, instruction := getActionAndInstruction(urlParams[2])
		found, mainAction, level, actionType := checkActionValidity(config.Actions, mainAction)
		log.Info("Checked instructions: <%s> <%s>", mainAction, instruction)
		if found == false {
			logs.Error(config.Elasticsearch.Url, config.Host, fmt.Sprintf("not found action <%s>, room <%s>, command <%s>", mainAction, concernedRoom, instruction))
			google_talk.Talk(ips, "Action introuvable")
			w.WriteHeader(500)
		} else {
			found := false
			if actionType == "domotiqueCommand" {
				logs.Info(config.Elasticsearch.Url, config.Host, fmt.Sprintf("Running action <%s>, room <%s>, command <%s>, level <%s>", mainAction, concernedRoom, instruction, level))
				found = RunDomoticCommand(config, instruction, concernedRoom, mainAction, level)
			}
			if found {
				w.WriteHeader(200)
			} else {
				logs.Error(config.Elasticsearch.Url, config.Host, fmt.Sprintf("not found action <%s>, room <%s>, command <%s>", mainAction, concernedRoom, instruction))
				google_talk.Talk(ips, "Instruction introuvable")
				w.WriteHeader(500)
			}
		}
	}
}

func extractDeviceDetails(config *configuration.Configuration, SearchedDevice string) (string, string, string, string) {
	for _, device := range config.Devices {
		if device.Name == SearchedDevice {
			fmt.Printf("found device : %+v \n", device)
			return device.Url, device.Instance, device.CommandClass, device.Id
		}
	}
	return "", "", "", ""
}

func ExecuteAction(config *configuration.Configuration, level string, deviceName string) (hasError bool) {
	hasError = false
	url, instance, commandClass, id := extractDeviceDetails(config, deviceName)
	if url != "" {
		err := ExecuteRequest(url, id, instance, commandClass, level)
		if err != nil {
			hasError = true
		}
	} else {
		fmt.Printf("Unable to find device %s \n", deviceName)
		hasError = true
	}
	return hasError
}

func AnalyseQuestionRequest(w http.ResponseWriter, r *http.Request, urlParams []string, config *configuration.Configuration) {
	requestType := urlParams[2]
	instruction := utils.ConvertInstruction(urlParams[3])
	log.Info("instructions: <%s> : <%s>", requestType, instruction)
	found := false
	foundText := ""
	if requestType == "listCommands" {
		for _, command := range config.Commands {
			for _, word := range command.Words {
				if strings.Contains(word, utils.RemoveEndletter(instruction)) {
					found = true
					foundText += "Allume ou éteins " + word + ";"
					time.Sleep(3 * time.Second)
				}
			}
		}
		//	google_talk.Talk(ips, foundText)
	}
	if found == false {
		//	google_talk.Talk(ips, "Je ne trouve aucune instruction contenant le mot " + instruction + ".")
	}
}
