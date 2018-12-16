package main

import (
	"fmt"
	"github.com/evalphobia/google-home-client-go/googlehome"
	"github.com/op/go-logging"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type Details struct {
	Url string `yaml:"url,omitempty"`
	Id []string `yaml:"id,omitempty"`
	Value string `yaml:"value,omitempty"`
	Instance string `yaml:"instance,omitempty"`
	CommandClass string `yaml:"commandClass,omitempty"`
}

type Command struct {
	Words string `yaml:"words,omitempty"`
	Actions []Details `yaml:"actions,omitempty"`
}

type Configuration struct {
	Commands []Command `yaml:"command,omitempty"`
	Cli *googlehome.Client
}


var log = logging.MustGetLogger("default")


var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{color:reset} %{message}`,
)


func ExecuteRequest(url string, id string, instance string, commandClass string, level string) (err error){
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	postingUrl := "http://192.168.222." + url + ":8083/ZWaveAPI/Run/devices[" + id + "].instances[" + instance + "].commandClasses["+ commandClass +"].Set("+ level + ")"
	log.Info("Request posted : %s", postingUrl)

	_, err = client.Get(postingUrl)
	if err != nil {
		fmt.Printf("Failed to execute request %s \n", postingUrl, err)
		return err
	}
	return nil
}


func main() {
	config := readConfiguration()
	cli, err := googlehome.NewClientWithConfig(googlehome.Config{
		Hostname: "192.168.222.135",
		Lang:     "fr",
		Accent:   "FR",
	})
	if err != nil {
		panic(err)
	}
	config.Cli = cli
	// Speak text on Google Home.


	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backend)
	backendLeveled.SetLevel(logging.NOTICE, "")
	logging.SetBackend(backendLeveled, backendFormatter)
	log.Info("Application Starting")

	http.HandleFunc("/switched/", func (w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path
		log.Info("Request received %s", urlPath)
		urlParams := strings.Split(urlPath, "/")
		if len(urlParams) == 3 {
			log.Info("Request succeeded")
			AnalyseRequest(w, r, urlParams)
		} else {
			log.Info("Request failed")
			w.WriteHeader(500)
		}
	})
	http.HandleFunc("/switch/", func (w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path
		urlParams := strings.Split(urlPath, "/")
		log.Info("Request received %s / %d", urlPath, len(urlParams))

		if len(urlParams) == 4 {
			log.Info("Request succeeded")
			AnalyseAIRequest(w, r, urlParams, config)
		} else {
			log.Info("Request failed")
			w.WriteHeader(500)
		}
	})

	err = http.ListenAndServe(":9998" , nil)
	if err != nil {
		log.Errorf("error %+v", err)
	}
}


func AnalyseRequest(w http.ResponseWriter, r *http.Request, urlParams []string) {
	level := urlParams[1]
	cmd:= urlParams[2]
	actions := strings.Split(cmd,"|")
	hasError := false;
	for _, action := range actions {
		values := strings.Split(action, ",")
		ids := strings.Split(values[0],"+")
		url := values[1]
		instance := values[2]
		commandClass := values[3]
		for _, id := range ids {
			err := ExecuteRequest(url, id, instance, commandClass, level)
			if err != nil {
				hasError = true
			}
		}
	}
	if hasError {
		w.WriteHeader(500)
	} else {
		w.WriteHeader(200)
	}
}

func AnalyseAIRequest(w http.ResponseWriter, r *http.Request, urlParams []string, config Configuration) {
	level := urlParams[2]
	instruction := strings.Replace(urlParams[3], "<<", "", 1)
	instruction = strings.Replace(instruction, ">>", "", 1)
	instruction = strings.Trim(instruction, " ")
	log.Info("instructions: <%s> : <%s>", level, instruction)
	for _, listAction := range config.Commands {
		if listAction.Words == instruction {
			config.Cli.SetLang("en")
			config.Cli.Notify("Instruction found")
		} else {
			config.Cli.SetLang("en")
			config.Cli.Notify("Instruction not found")
		}
	}
	//if hasError {
	//	w.WriteHeader(500)
	//} else {
		w.WriteHeader(200)
	//}
}

func readConfiguration() (Configuration) {
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

	var config Configuration

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Configuration Loaded : %+v \n", config)
	}
	return config
}
