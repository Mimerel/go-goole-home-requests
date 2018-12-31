package main

import (
	"fmt"
	"github.com/op/go-logging"
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

func AnalyseRequest(w http.ResponseWriter, r *http.Request, urlParams []string, config configuration.Configuration) {
	concernedRoom := urlParams[1]
	ips := []string{}
	for _, google := range config.Googles {
		if google.Name == concernedRoom {
			ips = google.Ip
		}
	}
	if len(ips) == 0 {
		config.Cli.Notify("Désolé, Enceinte pièce non reconnue")
		w.WriteHeader(500)
	} else {
		instruction := utils.ConvertInstruction(urlParams[2])
		log.Info("instructions: <%s> ", instruction)
		mainAction := strings.Split(instruction, " ")[0]
		instruction = strings.Replace(instruction, mainAction, "", 1)
		instruction = strings.Trim(instruction, "")
		level := ""
		found := false
		for _, action := range config.Actions {
			for _, actionName := range action.Name {
				if actionName == mainAction {
					level = action.Value
					mainAction = action.Replacement
					found = true
				}
			}
		}
		log.Info("instructions: <%s> <%s>",mainAction, instruction)
		if found == false {
			google_talk.Talk(ips, "Action introuvable")
			w.WriteHeader(500)
		} else {
			found = false
			for _, ListInstructions := range config.Commands {
				for _, word := range ListInstructions.Words {
					if utils.CompareWords(word, instruction, config) &&
						utils.CompareRooms(ListInstructions.Rooms, concernedRoom) &&
						utils.CompareActions(ListInstructions.Actions, mainAction) {
						for _, instruction := range ListInstructions.Instructions {
							if instruction.Value == "" {
								ExecuteAction(level, instruction.Instance, instruction.CommandClass, instruction.Url, instruction.Ids)

							} else {
								ExecuteAction(instruction.Value, instruction.Instance, instruction.CommandClass, instruction.Url, instruction.Ids)
							}
						}
						found = true;
						break
					}
				}
			}
			if found {
				w.WriteHeader(200)
			} else {
				google_talk.Talk(ips, "Instruction introuvable")
				w.WriteHeader(500)
			}
		}
	}
}

func ExecuteAction(level string, instance string, commandClass string, url string, ids []string) (hasError bool) {
	hasError = false;
	for _, id := range ids {
		err := ExecuteRequest(url, id, instance, commandClass, level)
		if err != nil {
			hasError = true
		}
	}
	return hasError
}

func AnalyseQuestionRequest(w http.ResponseWriter, r *http.Request, urlParams []string, config configuration.Configuration) {
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
