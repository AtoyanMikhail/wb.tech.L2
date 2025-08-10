package model

// Command описывает одну стадию пайплайна
type Command struct {
	Args      []string
	InFile    string
	OutFile   string
	AppendOut bool
}
