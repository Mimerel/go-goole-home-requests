package main

import (
	"fmt"
	"github.com/op/go-logging"
	"net/http"
	"os"
	"strings"
	"time"
)

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
	cli, err := googlehome.NewClientWithConfig(googlehome.Config{
		Hostname: "192.168.222.135",
		Lang:     "fr",
		Accent:   "FR",
	})
	if err != nil {
		panic(err)
	}

	// Speak text on Google Home.
	cli.SetLang("en")
	// cli.Notify("Google Home requester started")

	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backend)
	backendLeveled.SetLevel(logging.NOTICE, "")
	logging.SetBackend(backendLeveled, backendFormatter)
	log.Info("Application Starting")

	http.HandleFunc("/switch/", func (w http.ResponseWriter, r *http.Request) {
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