package googleAssistant

import (
	"fmt"
	"go-goole-home-requests/google_talk"
	"go-goole-home-requests/models"
	"go-goole-home-requests/utils"
	"go-goole-home-requests/zwave"
	"net/http"
	"strings"
)

/**
Method that searches for the ip(s) concerned by a room.
When an instruction is used, it will always be linked to a room
 */
func findIpOfGoogleHome(config *models.Configuration, concernedRoom string) ([]string) {
	ips := []string{}
	for _, google := range config.GoogleAssistant.GoogleBoxes {
		if google.Name == concernedRoom {
			ips = append(ips, google.Ip)
		}
	}
	return ips
}

/**
Method that splits the instruction into an action and a instruction
 */
func getActionAndInstruction(config *models.Configuration, instruction string) (action string, newInstruction string) {
	instruction = utils.ConvertInstruction(instruction)
	fmt.Printf("instructions: <%s> \n", instruction)
	mainAction := strings.Split(instruction, " ")[0]
	instruction = strings.Replace(instruction, mainAction, "", 1)
	instruction = strings.Trim(instruction, " ")
	return mainAction, instruction
}

/**
Method that checks if the action demanded exists and retrieves the information linked to this action.
 */
func checkActionValidity(config *models.Configuration, mainAction string) (string) {
	found := ""
	for _, action := range config.GoogleAssistant.GoogleTranslatedActionTypes {
		fmt.Printf("recherche correspondance de <%s> avec action :%v\n", mainAction, action)
		if action.ActionWord == mainAction {
			found = action.Action
			fmt.Printf("found\n")
			break
		}
	}
	return found
}

/**
Method that searches throw the list of possible commands for the
command sent by google home.
It first tries to find the corresponding "sentence" in its database.
IF it is found, it will check if the action is autorized in that room
If so, it will execute the command
 */
func RunDomoticCommand(config *models.Configuration, instruction string, concernedRoom string, mainAction string) (bool) {
	found := false
	for _, word := range config.GoogleAssistant.GoogleWords {
		fmt.Printf("comparing words <%s> with <%s>\n", instruction, word.Words)
		if utils.CompareWords(config, word.Words, instruction) {
			fmt.Printf("Found word\n")
			for _, ListInstructions := range config.GoogleAssistant.GoogleTranslatedInstructions {
				fmt.Printf("comparing room <%s> with room <%s> || action <%s> with <%s> || id <%v> with <%v>\n ",
					strings.ToUpper(concernedRoom), strings.ToUpper(ListInstructions.Room),strings.ToUpper(ListInstructions.Type) , strings.ToUpper(mainAction) ,
						word.ActionNameId , ListInstructions.ActionNameId )

				if strings.ToUpper(ListInstructions.Room) == strings.ToUpper(concernedRoom) &&
					strings.ToUpper(ListInstructions.Type) == strings.ToUpper(mainAction) &&
					word.ActionNameId == ListInstructions.ActionNameId {
					fmt.Printf("Found instruction %v\n", ListInstructions)
					ExecuteAction(config, ListInstructions)
					found = true;
				}
			}
			break
		}
	}
	return found
}

func AnalyseRequest(w http.ResponseWriter, r *http.Request, urlParams []string, config *models.Configuration) {
	concernedRoom := urlParams[1]
	ips := findIpOfGoogleHome(config, concernedRoom)
	if len(ips) == 0 {
		fmt.Printf("No google home ips found")
		w.WriteHeader(500)
	} else {
		mainAction, instruction := getActionAndInstruction(config, urlParams[2])
		mainAction = checkActionValidity(config, mainAction)
		fmt.Printf("Checked instructions: <%s> <%s>", mainAction, instruction)
		if mainAction == "" {
			config.Logger.Error(config.Elasticsearch.Url, config.Host, "not found action <%s>, room <%s>, command <%s>", mainAction, concernedRoom, instruction)
			google_talk.Talk(config, ips, "Action introuvable")
			fmt.Printf("Action introuvable")
			w.WriteHeader(500)
		} else {
			found := false
			config.Logger.Info(config.Elasticsearch.Url, config.Host, "Running action <%s>, room <%s>, command <%v>, level <%v>", mainAction, concernedRoom, instruction)
			found = RunDomoticCommand(config, instruction, concernedRoom, mainAction)
			if found {
				w.WriteHeader(200)
			} else {
				config.Logger.Error(config.Elasticsearch.Url, config.Host, "not found action <%s>, room <%s>, command <%s>", mainAction, concernedRoom, instruction)
				google_talk.Talk(config, ips, "Instruction introuvable")
				fmt.Printf("Instruction introuvable")

				w.WriteHeader(500)
			}
		}
	}
}

func ExecuteAction(config *models.Configuration, instruction models.GoogleTranslatedInstruction) (hasError bool) {
	err := zwave.ExecuteRequest(config, instruction.ZwaveUrl, instruction.DeviceId, instruction.Instance, instruction.CommandClass, instruction.Value)
	if err != nil {
		return true
	}
	return false
}
