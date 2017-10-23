package lib1ssarp

/*
TODO Создаем на основании конфига  проерку ролей доступа через переданный токен
TODO Он должен передаватся в не явном виде
TODO Создаем хэш со временем жизни сессии
TODO Уозволяем убивать сессию
TODO Проверяем роли
 */

import (
	"log"
	"fmt"
	"net/http"
	"regexp"
	"encoding/json"
	"bytes"
	"text/template"
	"io/ioutil"
	"strings"
	"crypto/md5"
	"sync"
)

const (
	//methods
	METHOD_GET = "GET"
	METHOD_POST = "POST"
	METHOD_PUT = "PUT"
	METHOD_DELETE = "DELETE"

	//rest api
	API_PREFIX = "api"

	//static page
	CLIENT_TEST = "/test-client"
	CLIENT_SESSION = "/session"
)

var regUriAll, regUriOne *regexp.Regexp

var synchronize = &sync.Mutex{}


func init()  {
	var e error
	regUriAll, e = regexp.Compile(`/` + API_PREFIX + `/([a-zA-Z]+)/?([a-zA-Z/]+)?$`)
	if e != nil {
		panic(e)
	}

	regUriOne, e = regexp.Compile(`/` + API_PREFIX + `/([a-zA-Z]+)/([0-9]+)/?([a-zA-Z/]+)?$`)
	if e != nil {
		panic(e)
	}
}


//Server
type HttpServer struct {
	Config Configuration
}



func (s HttpServer) Launch() {
	listen := s.Config.Server.Address()
	log.Println("listen: ", listen)

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc( s.index ))

	log.Fatal(http.ListenAndServe(listen, mux))
}


func (s HttpServer) index(w http.ResponseWriter, r *http.Request) {
	log.Println("Index page", r.URL.Path, r.Method);

	//simply route
	switch r.Method {
	case METHOD_GET:
		if r.URL.Path == CLIENT_TEST {
			s.clientTest(w, r)
		} else if regUriOne.MatchString(r.URL.Path) {
			s.one(w, r)
		} else {
			s.all(w, r)
		}
	case METHOD_POST:
		if r.URL.Path == CLIENT_SESSION {
			s.sessionOpen(w, r)
		} else {
			s.create(w, r)
		}
	case METHOD_PUT:
		s.update(w, r)
	case METHOD_DELETE:

		if r.URL.Path == CLIENT_SESSION {
			s.sessionClose(w, r)
		} else {
			s.delete(w, r)
		}
	}
}


//REST methods
//see http://www.restapitutorial.com/lessons/httpmethods.html
func (s HttpServer) all(w http.ResponseWriter, r *http.Request) {

	req := parseRequest(r.URL.Path , regUriAll , true)
	log.Println("Get all: ", req)

	m := s.findModel(req.Model)
	if m.Name != req.Model {
		status404(w)
		return
	}

	if !checkSession(m, ACTION_READ, r) {
		status401(w)
		return
	}

	ser := Service{s.Config.Database, m}
	d := ser.FetchAll()

	if d == nil {
		status404(w)
		return
	}

	js, e := json.Marshal(d)
	if e != nil {
		status500(w)
	}

	renderJson(w, string(js))
}

func (s HttpServer) one(w http.ResponseWriter, r *http.Request) {
	req := parseRequest(r.URL.Path , regUriOne , false)
	log.Println("Get model: ", req)

	if req.Id == "" {
		status404(w)
		return
	}

	m := s.findModel(req.Model)
	if m.Name != req.Model {
		status404(w)
		return
	}

	if !checkSession(m, ACTION_READ, r) {
		status401(w)
		return
	}

	ser := Service{s.Config.Database, m}
	d := ser.FetchOne(req.Id)

	if d == nil {
		status404(w)
		return
	}

	js, e := json.Marshal(d)
	if e != nil {
		log.Println(e)
		status500(w)
		return
	}

	renderJson(w, string(js))

}

func (s HttpServer) create(w http.ResponseWriter, r *http.Request) {
	req := parseRequest(r.URL.Path , regUriAll , false)
	log.Println("Create Model Name: ", req)

	body, e := ioutil.ReadAll(r.Body)
	if e != nil {
		log.Println(e)
		status500(w)
		return
	}

	log.Println("Body: ", string(body))

	var mt interface{}
	e = json.Unmarshal(body, &mt)
	if e != nil {
		log.Println(e)
		status500(w)
		return
	}

	mp := mt.(map[string]interface{})
	log.Println(mp)

	m := s.findModel(req.Model)
	if m.Name != req.Model {
		status404(w)
		return
	}

	if !checkSession(m, ACTION_CREATE, r) {
		status401(w)
		return
	}

	ser := Service{s.Config.Database, m}
	d := ser.Create(mp)

	log.Println("Create result: ", d)
	renderJson(w, fmt.Sprintf(`{"id": %d}`, d))
}

/**

 */
func (s HttpServer) update(w http.ResponseWriter, r *http.Request) {
	req := parseRequest(r.URL.Path , regUriOne , false)
	log.Println("Update Model Name: ", req)

	if req.Id == "" {
		status404(w)
		return
	}

	body, e := ioutil.ReadAll(r.Body)
	if e != nil {
		log.Println(e)
		status500(w)
		return
	}

	log.Println("Body: ", string(body))

	var mt interface{}
	e = json.Unmarshal(body, &mt)
	if e != nil {
		log.Println(e)
		status500(w)
		return
	}

	mp := mt.(map[string]interface{})
	log.Println(mp)

	m := s.findModel(req.Model)
	if m.Name != req.Model {
		status404(w)
		return
	}

	if !checkSession(m,ACTION_EDIT, r) {
		status401(w)
		return
	}

	ser := Service{s.Config.Database, m}
	d := ser.Update(req.Id, mp)

	log.Println("Update result: ", d)
	if d {
		renderJson(w, fmt.Sprintf(`{"update": "Ok"}`))
	} else {
		renderJson(w, fmt.Sprintf(`{"update": "No"}`))
	}
}

/**

 */
func (s HttpServer) delete(w http.ResponseWriter, r *http.Request) {
	req := parseRequest(r.URL.Path , regUriOne , false)
	log.Println("Delete Model Name: ", req)

	if req.Id == "" {
		status404(w)
		return
	}

	m := s.findModel(req.Model)
	if m.Name != req.Model {
		status404(w)
		return
	}

	if !checkSession(m, ACTION_DELETE, r) {
		status401(w)
		return
	}

	ser := Service{s.Config.Database, m}
	d := ser.Delete(req.Id)

	log.Println("Delete result: ", d)
	if d {
		renderJson(w, fmt.Sprintf(`{"delete": "Ok"}`))
	} else {
		renderJson(w, fmt.Sprintf(`{"delete": "No"}`))
	}
}


func (s HttpServer) sessionOpen(w http.ResponseWriter, r *http.Request) {

	synchronize.Lock()
	defer func() { synchronize.Unlock() }()

	log.Println("Open session... ")

	body, e := ioutil.ReadAll(r.Body)
	if e != nil {
		log.Println(e)
		status500(w)
		return
	}

	log.Println("Body: ", string(body))

	var data struct{
		Token string
		Sold string
	}

	e = json.Unmarshal(body, &data)
	if e != nil {
		log.Println(e)
		status500(w)
		return
	}

	//etch all allow token and check with hash
	for _, token := range s.Config.Session.Tokens {
		hash := digestString(token.Token + data.Sold + "POST")
		//log.Println(" hash: ", hash, "=", data.Token, "P:", token.Token + data.Sold + "POST")
		if hash == data.Token {
			log.Println("Find token: ", token)

			rl := make([]*Role, len(token.Roles) )
			//fill roles
			for i, rn := range token.Roles  {
				_, role := s.Config.Session.getRole(rn);
				rl[i] = role
			}

			hs := NewHttpSession( &token, rl )
			sessions[hs.CreatePublicToken()] = hs

			renderJson(w, fmt.Sprintf(`{"status": "OK", "session": "%s", "create": "%d"}`,
				hs.CreatePublicToken(),
				hs.Create))
			return
		}
	}

	//token not found
	status404(w)

}

func (s HttpServer) sessionClose(w http.ResponseWriter, r *http.Request) {

	synchronize.Lock()
	defer func() { synchronize.Unlock() }()

	log.Println("Close session: ")
}

/**

 */
func (s HttpServer) findModel(name string) Model {
	for _, m := range s.Config.Models {
		if m.Name == name {
			return m
		}
	}
	return Model{}
}

/**
Test page
 */
func (s HttpServer)  clientTest(w http.ResponseWriter, r *http.Request) {
	sourceFile := "./resource/client.html"
	fmt.Println("Read File: ", sourceFile)

	t, e := template.ParseFiles(sourceFile)
	if e != nil {
		log.Printf("Error Parse Template : %v\n", e)
		status404(w)
		return
	}

	var tpl bytes.Buffer

	e = t.Execute(&tpl, s.Config)
	if e != nil {
		log.Printf("Error Execute Template : %v\n", e)
		status500(w)
		return
	}

	fmt.Fprint(w, tpl.String())
}
//Server End


//helpers

type request struct {
	Model string
	Relation string
	Id string
}

func (r request) String() string {
	return fmt.Sprintf("request{Model: %s, Id: %s, Relation: %s}", r.Model, r.Id, r.Relation)
}

func renderJson(w http.ResponseWriter, js string) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(js))
}


func parseRequest(path string, req *regexp.Regexp, all bool ) request {
	res := req.FindStringSubmatch( path )

	var model, id, relation string

	switch len(res) {
	case 2:
		model = res[1]
	case 3:
		model = res[1]
		if all {
			relation = res[2]
		} else {
			id = res[2]
		}
	case 4:
		model = res[1]
		if all {
			relation = res[2]
		} else {
			id = res[2]
			relation = res[3]
		}
	}

	return request{Model:model, Id:id, Relation: relation}
}

/**
//TODO check session by header
 */
func checkSession (model Model, action ActionModel, r *http.Request) bool {

	synchronize.Lock()
	defer func() { synchronize.Unlock() }()

	token := r.Header.Get("Authorization")
	values := strings.Split(token, " ")

	if len(values) != 2 { 
		return false
	}

	switch strings.ToLower(values[0]) {
	case "token":
		log.Println("Authorization check token: ", values[1])

		if s, ok :=  sessions[values[1]]; ok {
			s.Update()
			log.Println("Find: ", s)
			return s.IsOpen() && s.CheckAccessModel(model, action)
		}
	}

	return false
}

func status404(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "custom 404")
}

func status500(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, "custom 500")
}

func status401(w http.ResponseWriter)  {
	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprint(w, "custom 401")
}

func digestString(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}