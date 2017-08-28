package lib1ssarp

import (
	"fmt"
)

func init() {}


type Configuration struct {
	Language string `language`
	Version string
	Output string
	Server Server
}

func (c Configuration) String () string {
	return fmt.Sprintf("Configuration{Language: %s, Version: %s, Server: %s}", c.Language, c.Version, c.Server)
}

type Server struct {
	Port uint
}

func (s Server) String () string {
	return fmt.Sprintf("Server{Port: %d}", s.Port)
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