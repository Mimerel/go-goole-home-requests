package configuration

import "github.com/evalphobia/google-home-client-go/googlehome"

type Details struct {
	Url          string   `yaml:"url,omitempty"`
	Ids          []string `yaml:"ids,omitempty"`
	Value        string   `yaml:"value,omitempty"`
	Instance     string   `yaml:"instance,omitempty"`
	CommandClass string   `yaml:"commandClass,omitempty"`
}

type Command struct {
	Words   []string  `yaml:"words,omitempty"`
	Rooms []string `yaml:"rooms,omitempty"`
	Instructions []Details `yaml:"instructions,omitempty"`
	Actions []string `yaml:"actions,omitempty"`
}

type Configuration struct {
	Commands []Command `yaml:"command,omitempty"`
	Cli      *googlehome.Client
	CharsToRemove []string `yaml:"charsToRemove,omitempty"`
	Googles []GoogleDetails `yaml:"googles,omitempty"`
	Actions []ActionDetails `yaml:"actions,omitempty"`
}

type GoogleDetails struct {
	Name string `yaml:"name,omitempty"`
	Ip []string `yaml:"ip,omitempty"`
}

type ActionDetails struct {
	Name []string `yaml:"name,omitempty"`
	Replacement string `yaml:"replacement,omitempty"`
	Value string `yaml:"value,omitempty"`
}
