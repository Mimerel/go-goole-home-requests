package configuration

import (
	"fmt"
	"go-goole-home-requests/models"
	"go-goole-home-requests/utils"
	"github.com/Mimerel/go-utils"
)

func executeGoogleAssistantConfiguration(config *models.Configuration) {
	fmt.Printf("Collecting Database Data\n")
	err := getDevices(config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Collecting Box information\n")
	err = getBoxes(config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Collecting Google Words information\n")
	err = getWords(config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Collecting Google Instructions\n")
	err = getGoogleInstructions(config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Collecting Google Action Names\n")
	err = getActionNames(config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Collecting Rooms\n")
	err = getRooms(config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Collecting Device Types\n")
	err = getDeviceTypes(config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Collecting Google Boxes\n")
	err = getGoogleBox(config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Collecting Google Boxes\n")
	err = getGoogleActionTypes(config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Collecting Google Boxes\n")
	err = getGoogleActionTypesWords(config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Checking configuration\n")
	checkConfiguration(config)
	fmt.Printf("saving translated instructions\n")
	saveConfigToDataBase(config)
}

func saveConfigToDataBase(c *models.Configuration) {
	db := utils.CreateDbConnection(c)
	db.Debug = false
	fmt.Printf("Emptied instructions\n")
	db.Request("delete from " + utils.TableGoogleTranslatedInstructions)
	fmt.Printf("saving instructions\n")
	for k,v := range c.GoogleAssistant.GoogleTranslatedInstructions {
		fmt.Printf("key : %v => %v\n", k, v)
	}
	err := utils.ActionInMariaDB(c, c.GoogleAssistant.GoogleTranslatedInstructions, utils.TableGoogleTranslatedInstructions, utils.ActionInsertIgnore)
	if err != nil {
		c.Logger.Error("Unable to store request model in MariaDB : %+v", err)
	}
}

func getBoxes(c *models.Configuration) (err error) {
	db := utils.CreateDbConnection(c)
	db.Table = utils.TableDomotiqueBox
	db.WhereClause = ""
	db.Debug = false
	db.DataType = new([]models.Zwave)
	res, err := go_utils.SearchInTable(db)
	if err != nil {
		c.Logger.Error("Unable to request database : %v", err)
		return err
	}
	if len(*res.(*[]models.Zwave)) > 0 {
		c.Zwaves = *res.(*[]models.Zwave)
		return nil
	}
	return fmt.Errorf("Unable to find list of Zwave Boxes")
}

func getGoogleActionTypes(c *models.Configuration) (err error) {
	db := utils.CreateDbConnection(c)
	db.Table = utils.TableGoogleActionTypes
	db.WhereClause = ""
	db.Debug = false
	db.DataType = new([]models.GoogleActionTypes)
	res, err := go_utils.SearchInTable(db)
	if err != nil {
		c.Logger.Error("Unable to request database : %v", err)
		return err
	}
	if len(*res.(*[]models.GoogleActionTypes)) > 0 {
		c.GoogleAssistant.GoogleActionTypes = *res.(*[]models.GoogleActionTypes)
		return nil
	}
	return fmt.Errorf("Unable to find list of Zwave Boxes")
}

func getGoogleActionTypesWords(c *models.Configuration) (err error) {
	db := utils.CreateDbConnection(c)
	db.Table = utils.TableGoogleActionTypesWords
	db.WhereClause = ""
	db.Debug = false
	db.DataType = new([]models.GoogleActionTypesWords)
	res, err := go_utils.SearchInTable(db)
	if err != nil {
		c.Logger.Error("Unable to request database : %v", err)
		return err
	}
	if len(*res.(*[]models.GoogleActionTypesWords)) > 0 {
		c.GoogleAssistant.GoogleActionTypesWords = *res.(*[]models.GoogleActionTypesWords)
		return nil
	}
	return fmt.Errorf("Unable to find list of Zwave Boxes")
}

func getDevices(c *models.Configuration) (err error) {
	db := utils.CreateDbConnection(c)
	db.Table = utils.TableDevices
	db.WhereClause = ""
	db.Seperator = ","
	db.Debug = false
	db.DataType = new([]models.DeviceDetails)
	res, err := go_utils.SearchInTable(db)
	if err != nil {
		c.Logger.Error("Unable to request database for devices: %v", err)
		return err
	}
	if len(*res.(*[]models.DeviceDetails)) > 0 {
		c.Devices = *res.(*[]models.DeviceDetails)
		return nil
	}
	return fmt.Errorf("Unable to find list of Devices")
}

func getWords(c *models.Configuration) (err error) {
	db := utils.CreateDbConnection(c)
	db.Table = utils.TableGoogleWords
	db.WhereClause = ""
	db.Seperator = ","
	db.Debug = false
	db.DataType = new([]models.GoogleWords)
	res, err := go_utils.SearchInTable(db)
	if err != nil {
		c.Logger.Error("Unable to request database for words: %v", err)
		return err
	}
	if len(*res.(*[]models.GoogleWords)) > 0 {
		c.GoogleAssistant.GoogleWords = *res.(*[]models.GoogleWords)
		return nil
	}
	return fmt.Errorf("Unable to find list of words")
}

func getGoogleInstructions(c *models.Configuration) (err error) {
	db := utils.CreateDbConnection(c)
	db.Table = utils.TableGoogleInstructions
	db.WhereClause = ""
	db.Seperator = ","
	db.Debug = false
	db.DataType = new([]models.GoogleInstruction)
	res, err := go_utils.SearchInTable(db)
	if err != nil {
		c.Logger.Error("Unable to request database for words: %v", err)
		return err
	}
	if len(*res.(*[]models.GoogleInstruction)) > 0 {
		c.GoogleAssistant.GoogleInstructions = *res.(*[]models.GoogleInstruction)
		return nil
	}
	return fmt.Errorf("Unable to find list of words")
}

func getActionNames(c *models.Configuration) (err error) {
	db := utils.CreateDbConnection(c)
	db.Table = utils.TableGoogleActionNames
	db.WhereClause = ""
	db.Seperator = ","
	db.Debug = false
	db.DataType = new([]models.GoogleActionNames)
	res, err := go_utils.SearchInTable(db)
	if err != nil {
		c.Logger.Error("Unable to request database for words: %v", err)
		return err
	}
	if len(*res.(*[]models.GoogleActionNames)) > 0 {
		c.GoogleAssistant.GoogleActionNames = *res.(*[]models.GoogleActionNames)
		return nil
	}
	return fmt.Errorf("Unable to find list of words")
}

func getRooms(c *models.Configuration) (err error) {
	db := utils.CreateDbConnection(c)
	db.Table = utils.TableRooms
	db.WhereClause = ""
	db.Seperator = ","
	db.Debug = false
	db.DataType = new([]models.Room)
	res, err := go_utils.SearchInTable(db)
	if err != nil {
		c.Logger.Error("Unable to request database for words: %v", err)
		return err
	}
	if len(*res.(*[]models.Room)) > 0 {
		c.Rooms = *res.(*[]models.Room)
		return nil
	}
	return fmt.Errorf("Unable to find list of words")
}

func getGoogleBox(c *models.Configuration) (err error) {
	db := utils.CreateDbConnection(c)
	db.Table = utils.TableGoogleBox
	db.WhereClause = ""
	db.Seperator = ","
	db.Debug = false
	db.DataType = new([]models.GoogleBox)
	res, err := go_utils.SearchInTable(db)
	if err != nil {
		c.Logger.Error("Unable to request database for words: %v", err)
		return err
	}
	if len(*res.(*[]models.GoogleBox)) > 0 {
		c.GoogleAssistant.GoogleBoxes = *res.(*[]models.GoogleBox)
		return nil
	}
	return fmt.Errorf("Unable to find list of words")
}

func getDeviceTypes(c *models.Configuration) (err error) {
	db := utils.CreateDbConnection(c)
	db.Table = utils.TableDeviceTypes
	db.WhereClause = ""
	db.Seperator = ","
	db.Debug = false
	db.DataType = new([]models.DeviceType)
	res, err := go_utils.SearchInTable(db)
	if err != nil {
		c.Logger.Error("Unable to request database for words: %v", err)
		return err
	}
	if len(*res.(*[]models.DeviceType)) > 0 {
		c.DeviceTypes = *res.(*[]models.DeviceType)
		return nil
	}
	return fmt.Errorf("Unable to find list of words")
}

/**
Method that checks that the configuration file is consistent.
If a device name that does not exist in the device list
or if a zwave device that does not exist in the zwave list
are used, error message will be displayed and the program will stop
 */
func checkConfiguration(config *models.Configuration) {
	// Check if devices that are used in commands are in the device list
	//[]GoogleTranslatedInstruction

	for _, device := range config.Devices {
		translated := new(models.DeviceTranslated)
		translated.DomotiqueId = device.DomotiqueId
		translated.Instance = device.Instance
		translated.CommandClass = device.CommandClass
		translated.DeviceId = device.DeviceId
		translated.Room = getRoomFromId(config, device.RoomId).Name
		translated.Type = getTypeFromId(config, device.TypeId).Name
		translated.Name = device.Name
		translated.ZwaveName = getZwaveFromId(config, device.Zwave).Name
		translated.ZwaveUrl = getZwaveFromId(config, device.Zwave).Ip
		config.DevicesTranslated = append(config.DevicesTranslated, *translated)
	}

	for _, v := range config.GoogleAssistant.GoogleActionTypesWords {
		translated := new(models.GoogleTranslatedActionTypes)
		translated.ActionWord = v.Action
		translated.Action = getActionTypeFromId(config, v.ActionTypeId).Name
		config.GoogleAssistant.GoogleTranslatedActionTypes = append(config.GoogleAssistant.GoogleTranslatedActionTypes, *translated)
	}

	for _, instruction := range config.GoogleAssistant.GoogleInstructions {
		translated := new(models.GoogleTranslatedInstruction)
		translated.Value = instruction.Value
		translated.Id = instruction.Id
		translated.ActionName = getActionNameFromId(config, instruction.ActionNameId).Name
		translated.ActionNameId = instruction.ActionNameId
		translated.CommandClass = getDeviceFromId(config, instruction.DomotiqueId).CommandClass
		translated.Instance = getDeviceFromId(config, instruction.DomotiqueId).Instance
		translated.DeviceName = getDeviceFromId(config, instruction.DomotiqueId).Name
		translated.DeviceId =  getDeviceFromId(config, instruction.DomotiqueId).DeviceId
		translated.Room = getDeviceFromId(config, instruction.DomotiqueId).Room
		translated.ZwaveUrl = getDeviceFromId(config, instruction.DomotiqueId).ZwaveUrl
		translated.TypeDevice = getDeviceFromId(config, instruction.DomotiqueId).Type
		translated.Type = getActionTypeFromId(config, instruction.TypeId).Name
		config.GoogleAssistant.GoogleTranslatedInstructions = append(config.GoogleAssistant.GoogleTranslatedInstructions, *translated)
	}

}

func getRoomFromId(c *models.Configuration, id int64) (models.Room) {
	for _, v := range c.Rooms {
		if v.Id == id {
			return v
		}
	}
	return models.Room{}
}

func getZwaveFromId(c *models.Configuration, id int64) (models.Zwave) {
	for _, v := range c.Zwaves {
		if v.Id == id {
			return v
		}
	}
	return models.Zwave{}
}

func getTypeFromId(c *models.Configuration, id int64) (models.DeviceType) {
	for _, v := range c.DeviceTypes {
		if v.Id == id {
			return v
		}
	}
	return models.DeviceType{}
}

func getActionNameFromId(c *models.Configuration, id int64) (models.GoogleActionNames) {
	for _, v := range c.GoogleAssistant.GoogleActionNames {
		if v.Id == id {
			return v
		}
	}
	return models.GoogleActionNames{}
}

func getDeviceFromId(c *models.Configuration, id int64) (models.DeviceTranslated) {
	for _, v := range c.DevicesTranslated {
		if v.DomotiqueId == id {
			return v
		}
	}
	return models.DeviceTranslated{}
}

func getActionTypeFromId(c *models.Configuration, id int64) (models.GoogleActionTypes) {
	for _, v := range c.GoogleAssistant.GoogleActionTypes {
		if v.Id == id {
			return v
		}
	}
	return models.GoogleActionTypes{}
}

