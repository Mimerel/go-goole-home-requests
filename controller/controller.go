package controller

import (
	"fmt"
	"go-goole-home-requests/configuration"
	"net/http"
)

func Controller() {
	config := configuration.ReadConfiguration()

	config.Logger.Info("Application Starting")

	getControllerGoogleAssistant(config)

	err := http.ListenAndServe(":9998", nil)
	if err != nil {
		fmt.Printf("error %+v", err)
	}
}