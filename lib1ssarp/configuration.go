package lib1ssarp

import (
	"fmt"
)

func init() {}


type Configuration struct {
	Language string `language`
	Version string
	Output string
}

func (c Configuration) String () string {
	return fmt.Sprintf("Configuration{Language: %s, Version: %s}", c.Language, c.Version)
}

type Models struct {
	models []Model
}

type Model struct {
	Name string
	fields []Field
}

type Field struct {
	Name string
	Type string
	Autoincrement bool
	Length uint
}


type ConfigurationBuilder interface {
	Build(c Configuration)
}