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

//Server
type Server struct {
	Port uint
	Host string
}

func (s Server) String () string {
	return fmt.Sprintf("Server{Address: %s}", s.Address())
}

func (s Server) Address () string {
	h := s.Host
	if h == "" {
		h = "localhost"
	}
	p := s.Port
	if p == 0 {
		p = 8090
	}
	return fmt.Sprintf("%s:%d", h, p)
}

//Server end

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
	RunServer(c Configuration)
}