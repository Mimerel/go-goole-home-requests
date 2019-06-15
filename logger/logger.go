package logger

import (
	"fmt"
	"go-goole-home-requests/models"
	"go-goole-home-requests/utils"
	"strings"
	"time"
)

func Info(config *models.Configuration, module string, message string, args ...interface{}) {
	computedMessage := fmt.Sprintf(message, args...)
	fmt.Printf(time.Now().Format(time.RFC3339)+" - Info (%s): %s \n", module, computedMessage)
	sendLogToDB(config, module, computedMessage)
}

func Debug(config *models.Configuration, module string, message string, args ...interface{}) {
	computedMessage := fmt.Sprintf(message, args...)
	fmt.Printf(time.Now().Format(time.RFC3339)+" - Debug (%s): %s \n", module, computedMessage)
	sendLogToDB(config, module, computedMessage)
}

func Error(config *models.Configuration, module string, message string, args ...interface{}) {
	computedMessage := fmt.Sprintf(message, args...)
	fmt.Printf(time.Now().Format(time.RFC3339)+" - Error (%s): %s \n", module, computedMessage)
	sendLogToDB(config, module, computedMessage)
}

func sendLogToDB(c *models.Configuration, module string, computedMessage string) {
	db := utils.CreateDbConnection(c)
	db.Database = utils.LoggerDatabase
	db.Debug = true
	logs := []models.Log{models.Log{Module: module, Message: computedMessage}}

	col, val, err := db.DecryptStructureAndData(logs)
	if err != nil {
		c.Logger.Error("col %s", col)
		c.Logger.Error("val %s", val)
	}
	err = db.Insert(false, utils.LoggerDomotique, col, val)

	if err != nil {
		c.Logger.Error("table %s", utils.LoggerDomotique)
		c.Logger.Error("col %s", col)
		values := strings.Split(val, "),(")
		for k, v := range values {
			c.Logger.Error("row %v - %s", k, v)
		}
	}
}
