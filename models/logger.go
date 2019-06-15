package models

type Log struct {
	Module   string `csv:"module"`
	Message  string `csv:"message"`
}
