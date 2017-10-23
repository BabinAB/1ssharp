package lib1ssarp

import (
	"time"
	"crypto/sha256"
	"fmt"
)

type ActionModel int

const (
	TIME_LIVE int64 = 86400
	SECRET = "123fefefvvvf#222121+87^323lcscc11122njncjc097663th((csacsacsacs"
	ACTION_READ ActionModel = 1
	ACTION_EDIT ActionModel = 2
	ACTION_DELETE ActionModel = 4
	ACTION_CREATE ActionModel = 8

)


var sessions map[string]*HttpSession

func init()  {
	sessions = make(map[string] *HttpSession)
}


func NewHttpSession(t *Token, r []*Role) *HttpSession {
	unix := time.Now().Unix()
	s := HttpSession{Token: t, Roles: r, open: true, LastUpdate: unix, Create: unix}
	s.Update()
	return &s
}


type HttpSession struct {
	Token *Token
	Roles []*Role
	LastUpdate int64
	Create int64
	open bool
}

func (s *HttpSession) String()string  {
	return fmt.Sprintf("HttpSession{Open: %x, Create: %d, LastUpdate: %d}",
		s.IsOpen(), s.Create, s.LastUpdate)
}

func (s *HttpSession) Close()  {
	s.open = false
}

func (s *HttpSession) Open()  {
	fmt.Println(s.open)
	s.open = true
	fmt.Println(s.open)
}

func (s *HttpSession) IsOpen() bool  {
	return s.open
}

func (s *HttpSession) Update() {

	if !s.IsOpen() {
		return
	}

	t := time.Now().Unix()

	if s.LastUpdate + TIME_LIVE > t {
		s.LastUpdate = t
	} else {
		s.Close()
	}
}

/**

 */
func (s *HttpSession) CreatePublicToken() string {
	str := s.Token.Token + SECRET + fmt.Sprintf("%d", s.Create)
	return fmt.Sprintf("%x", sha256.Sum256([]byte(str)))
}

//TODO check access model name
func (s *HttpSession) CheckAccessModel(m Model, action ActionModel) bool {
	for _, r := range s.Roles {
		perm := r.getPermission(m.Name)
		if nil != perm && perm.CheckAction(action) {
			return true
		}
	}

	return false
}


