package models

type Elasticsearch struct {
	Url string `yaml:"url,omitempty"`
}

type GoogleWords struct {
	Words        string `csv:"words"`
	Id           int64  `csv:"id"`
	ActionNameId int64  `csv:"actionNameId"`
}

type GoogleActionNames struct {
	Id   int64  `csv:"id"`
	Name string `csv:"name"`
}

type GoogleInstruction struct {
	Id           int64 `csv:"id"`
	ActionNameId int64 `csv:"actionNameId"`
	TypeId       int64 `csv:"typeId"`
	DomotiqueId  int64 `csv:"domotiqueId"`
	GoogleBoxId  int64 `csv:"googleBoxId"`
	Value        int64 `csv:"value"`
}

type GoogleTranslatedInstruction struct {
	Id           int64  `csv:"id"`
	ActionName   string `csv:"actionName"`
	ActionNameId int64  `csv:"actionNameId"`
	Type         string `csv:"type"`
	DeviceName   string `csv:"deviceName"`
	DeviceId     int64  `csv:"deviceId"`
	ZwaveUrl     string `csv:"zwaveUrl"`
	Room         string `csv:"room"`
	TypeDevice   string `csv:"typeDevice"`
	Instance     int64  `csv:"instance"`
	CommandClass int64  `csv:"commandClass"`
	Value        int64  `csv:"value"`
}

type DeviceDetails struct {
	Name         string `csv:"name"`
	DeviceId     int64  `csv:"deviceId"`
	DomotiqueId  int64  `csv:"domotiqueId"`
	RoomId       int64  `csv:"roomId"`
	TypeId       int64  `csv:"typeId"`
	Zwave        int64  `csv:"boxId"`
	Instance     int64  `csv:"instance"`
	CommandClass int64  `csv:"commandClass"`
}

type DeviceTranslated struct {
	Name         string
	DeviceId     int64
	DomotiqueId  int64
	Room         string
	Type         string
	ZwaveName    string
	ZwaveUrl     string
	Instance     int64
	CommandClass int64
}

type Zwave struct {
	Id   int64  `csv:"id"`
	Name string `csv:"name"`
	Ip   string `csv:"ip"`
}

type GoogleActionTypes struct {
	Id   int64  `csv:"id"`
	Name string `csv:"name"`
}

type GoogleActionTypesWords struct {
	Id           int64  `csv:"id"`
	ActionTypeId int64  `csv:"remplacementId"`
	Action       string `csv:"name"`
}

type GoogleTranslatedActionTypes struct {
	ActionWord string
	Action     string
}

type Room struct {
	Id   int64  `csv:"id"`
	Name string `csv:"room"`
}

type DeviceType struct {
	Id   int64  `csv:"id"`
	Name string `csv:"name"`
}

type GoogleBox struct {
	Id   int64  `csv:"id"`
	Name string `csv:"name"`
	Ip   string `csv:"ip"`
}

type MariaDB struct {
	User     string `yaml:"user,omitempty"`
	Password string `yaml:"password,omitempty"`
	IP       string `yaml:"ip,omitempty"`
	Port     string `yaml:"port,omitempty"`
	Database string `yaml:"database,omitempty"`
}
