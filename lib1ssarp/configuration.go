package lib1ssarp

import (
	"fmt"
)

func init() {}


type Configuration struct {
	Language string `language`
	Version string
	Output string
	Database Database
	Server Server
	Models []Model
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


//Database
type Database struct {
	Type string `mysql` //mysql | sqlite
	Host string `localhost`
	Port uint `3306`
	Basename string
	Username string
	Password string
	Path string 	//for  sqlite ./db.sqlite
}

//Database end

//Model
type Model struct {
	Name string
	Fields []Field
}

func (m Model) String() string {
	return fmt.Sprintf("Model{Name: %s, FieldLen: %d}", m.Name, len(m.Fields))
}
//Model end


//Field
type Field struct {
	Title string //Title for view label
	Name string  //Name for DB field
	Type string
	Autoincrement bool
	Length uint //Length for string type
	Refer ReferField
}

func (f Field) String() string {
	return fmt.Sprintf("Field{Name: %s, Type: %s}", f.Name, f.Type)
}
//Field end


//Refer
type ReferField struct {
	Name string
}

func (r ReferField) String() string {
	return fmt.Sprintf("ReferField{Name: %s}",  r.Name)
}
//Refer end



type ConfigurationBuilder interface {
	Build(c Configuration)
	RunServer(c Configuration)
}