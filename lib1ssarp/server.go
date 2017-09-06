package lib1ssarp

/*
TODO move to main package
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
)

const METHOD_GET = "GET"
const METHOD_POST = "POST"
const METHOD_PUT = "PUT"
const METHOD_DELETE = "DELETE"

const API_PREFIX = "api"

const CLIENT_TEST = "/test-client"


var regUriAll *regexp.Regexp
var regUriOne *regexp.Regexp


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
		s.create(w, r)
	case METHOD_PUT:
		s.update(w, r)
	case METHOD_DELETE:
		s.delete(w, r)
	}
}


//REST methods
//see http://www.restapitutorial.com/lessons/httpmethods.html
func (s HttpServer) all(w http.ResponseWriter, r *http.Request) {

	res := regUriAll.FindStringSubmatch( r.URL.Path )

	var model, relation string

	switch len(res) {
	case 2:
		model = res[1]
	case 3:
		model = res[1]
		relation = res[2]
	default:
		status404(w)
		return
	}

	log.Println("Model Name: ", model, ",  relation: ", relation)

	m := s.findModel(model)
	if m.Name != model {
		status404(w)
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

	res := regUriOne.FindStringSubmatch( r.URL.Path )

	var model, id, relation string

	switch len(res) {
	case 3:
		model = res[1]
		id = res[2]
	case 4:
		model = res[1]
		id = res[2]
		relation = res[3]
	default:
		status404(w)
		return
	}

	log.Println("Model Name: ", model, ", pk: ", id, ", relation: ", relation)

	m := s.findModel(model)
	if m.Name != model {
		status404(w)
		return
	}

	ser := Service{s.Config.Database, m}
	d := ser.FetchOne(id)

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
	res := regUriAll.FindStringSubmatch( r.URL.Path )

	var model string

	switch len(res) {
	case 2:case 3:
		model = res[1]
	default:
		status404(w)
		return
	}

	log.Println("Create Model Name: ", model)

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

	m := s.findModel(model)
	if m.Name != model {
		status404(w)
		return
	}

	ser := Service{s.Config.Database, m}
	d := ser.Create(mp)

	log.Println("Create result: ", d)
	renderJson(w, fmt.Sprintf("{id: %d}", d))
}

/**

 */
func (s HttpServer) update(w http.ResponseWriter, r *http.Request) {

}

/**

 */
func (s HttpServer) delete(w http.ResponseWriter, r *http.Request) {

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

func renderJson(w http.ResponseWriter, js string) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(js))
}



func status404(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "custom 404")
}

func status500(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, "custom 500")
}