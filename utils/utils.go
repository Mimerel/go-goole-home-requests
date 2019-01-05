package utils

import (
	"strings"
)

func ConvertInstruction (value string) string {
	instruction := strings.Replace(value, "<<", "", 1)
	instruction = strings.Replace(instruction, ">>", "", 1)
	instruction = strings.Trim(instruction, " ")
	newValue := ""
	for i := 0; i < len(instruction); i++ {
		newValue = newValue + string(instruction[i])
	}
	return newValue
}

func RemoveEndletter( value string) string {
	if strings.HasSuffix(value, "s") {
		value = strings.TrimSuffix(value, "s")
	} else if strings.HasSuffix(value, "x"){
		value = strings.TrimSuffix(value, "x")
	}
	return value
}

func CompareWords(word string, instruction string ) (bool) {
	same := true;
	if strings.ToLower(strings.Replace(word, " ", "", -1)) != strings.ToLower(strings.Replace(instruction, " ", "", -1)) {
		same = false
	}
	return same
}

func CompareRooms(rooms []string, requestedFrom string) (bool) {
	same := false;
	for _, value := range rooms {
		if value == requestedFrom {
			same = true
		}
	}
	return same
}

func CompareActions(actions []string, RequestedAction string) (bool) {
	same := false;
	for _, value := range actions {
		if value == RequestedAction {
			same = true
		}
	}
	return same
}

func IsInArray(list []string, value string) (bool) {
	exists := false;
	for _, v := range list {
		if v == value {
			exists = true
		}
	}
	return exists
}