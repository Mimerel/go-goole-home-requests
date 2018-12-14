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
	postingUrl := "http://" + url + ":8083/ZWaveAPI/Run/devices[" + id + "].instances[" + instance + "].commandClasses["+ commandClass +"].Set("+ level + ")"
	log.Info("Request posted : %s", postingUrl)

	_, err = client.Get(postingUrl)
	if err != nil {
		fmt.Printf("Failed to execute request %s \n", postingUrl, err)
		return err
	}
	return nil
}


func main() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backend)
	backendLeveled.SetLevel(logging.NOTICE, "")
	logging.SetBackend(backendLeveled, backendFormatter)
	log.Info("Appliciation Starting")

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		log.Info("Request received")
		urlPath := r.URL.Path
		urlParams := strings.Split(urlPath, "/")
		if len(urlParams) == 3 {
			log.Info("Request succeeded")
			AnalyseRequest(w, r, urlParams)
		} else {
			log.Info("Request failed")

			w.WriteHeader(500)
		}
		})
	err := http.ListenAndServe(":9998" , nil)
	if err != nil {
		log.Errorf("error %+v", err)
	}
}


func AnalyseRequest(w http.ResponseWriter, r *http.Request, urlParams []string) {
	level := urlParams[1]
	cmd:= urlParams[2]
	values := strings.Split(cmd, ",")
	id := values[0]
	url := values[1]
	instance := values[2]
	commandClass := values[3]
	err := ExecuteRequest(url, id, instance, commandClass, level)
	if err != nil {
		w.WriteHeader(500)
	} else {
		w.WriteHeader(200)
	}
}