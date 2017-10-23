package lib1ssarp

import (
	"fmt"
	"errors"
)

func init() {}


type Configuration struct {
	Language string `language`
	Version string
	Output string
	Database Database
	Server Server
	Models []Model
	Session Session
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



//Session

type Permission struct {
	Model string
	Read bool
	Update bool
	Delete bool
	Create bool
}

func (p *Permission) String () string {
	return fmt.Sprintf("Permission{Name: %s}",  p.Model)
}


func (p *Permission) CheckAction(action ActionModel) bool {
	switch action {
	case ACTION_CREATE:
		return p.Create
	case ACTION_DELETE:
		return p.Delete
	case ACTION_EDIT:
		return p.Update
	case ACTION_READ:
		return p.Read
	}
	return false
}


type Role  struct {
	Name string
	Permissions []Permission
}

func (r *Role) String () string {
	return fmt.Sprintf("Role{Name: %s}",  r.Name)
}

/**
Get permission role by nameModel(name permission)
 */
func (r *Role) getPermission ( modelName string) *Permission {
	for _, p := range r.Permissions {
		if p.Model == modelName {
			return &p
		}
	}
	return nil
}

type Token struct {
	Token string
	Roles []string
}

func (t Token) String () string {
	return fmt.Sprintf("Token{Name: %s}",  t.Token)
}

type Session struct {
	Roles []Role
	Tokens []Token
}

func (s Session) getToken(token string) (error, Token) {

	for _, t := range s.Tokens {
		if t.Token == token {
			return nil, t
		}
	}

	return errors.New("Not found"), Token{}
}

func (s Session) getRole(name string) (error, *Role) {

	for _, r := range s.Roles {
		if r.Name == name {
			return nil, &r
		}
	}

	return  errors.New("Not found"), nil
}


//Session end

type ConfigurationBuilder interface {
	Build(c Configuration)
	RunServer(c Configuration)
}