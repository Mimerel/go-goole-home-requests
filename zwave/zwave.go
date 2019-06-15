package zwave

import (
	"go-goole-home-requests/logger"
	"net/http"
	"strconv"
	"time"
	"go-goole-home-requests/models"
)

/**
Method that sends a request to a domotic zwave server to run an instruction
 */
func ExecuteRequest(config *models.Configuration, url string, id int64, instance int64, commandClass int64, level int64) (err error) {
	logger.Info(config, "ExecuteRequest", "Préparing post")
	timeout := time.Duration(20 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	postingUrl := "http://" + url + ":8083/ZWaveAPI/Run/devices[" + strconv.FormatInt(id, 10) + "].instances[" + strconv.FormatInt(instance, 10) + "].commandClasses[" + strconv.FormatInt(commandClass, 10) + "].Set(" + strconv.FormatInt(level, 10) + ")"
	logger.Info(config, "ExecuteRequest", "Request posted : %s", postingUrl)

	_, err = client.Get(postingUrl)
	if err != nil {
		logger.Error(config, "ExecuteRequest", "Failed to execute request %s \n", postingUrl, err)
		return err
	}
	logger.Info(config, "ExecuteRequest", "Request successful...")
	return nil
}

