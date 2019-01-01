package configuration

type Details struct {
	Url          string   `yaml:"url,omitempty"`
	Ids          []string `yaml:"ids,omitempty"`
	Instance     string   `yaml:"instance,omitempty"`
	CommandClass string   `yaml:"commandClass,omitempty"`
	DeviceName string     `yaml:"name,omitempty"`
	Value        string   `yaml:"value,omitempty"`
}

type Command struct {
	Words   []string  `yaml:"words,omitempty"`
	Rooms []string `yaml:"rooms,omitempty"`
	Instructions []Details `yaml:"instructions,omitempty"`
	Actions []string `yaml:"actions,omitempty"`
}

type Configuration struct {
	Commands []Command `yaml:"command,omitempty"`
	CharsToRemove []string `yaml:"charsToRemove,omitempty"`
	Googles []GoogleDetails `yaml:"googles,omitempty"`
	Actions []ActionDetails `yaml:"actions,omitempty"`
	Devices []Device `yaml:"devices,omitempty"`
	Zwaves []Zwave `yaml:"zwaves,omitempty"`
}

type Device struct {
	Name         string `yaml:"name,omitempty"`
	Id           string `yaml:"id,omitempty"`
	Url          string   `yaml:"url,omitempty"`
	Zwave        string   `yaml:"zwave,omitempty"`
	Instance     string   `yaml:"instance,omitempty"`
	CommandClass string   `yaml:"commandClass,omitempty"`
}

type GoogleDetails struct {
	Name string `yaml:"name,omitempty"`
	Ip []string `yaml:"ip,omitempty"`
}

type Zwave struct {
	Name string `yaml:"name,omitempty"`
	Ip string `yaml:"ip,omitempty"`
}

type ActionDetails struct {
	Name []string `yaml:"name,omitempty"`
	Replacement string `yaml:"replacement,omitempty"`
	Type string `yaml:"type,omitempty"`
	Value string `yaml:"value,omitempty"`
}
