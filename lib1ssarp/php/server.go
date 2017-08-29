package lib1ssarp_php

/*
TODO move to main package
 */

import (
	"./../../lib1ssarp"
	"log"
	"fmt"
	"net/http"
	"regexp"
)

const METHOD_GET = "GET"
const METHOD_POST = "POST"
const METHOD_PUT = "PUT"
const METHOD_DELETE = "DELETE"


var regUriAll *regexp.Regexp
var regUriOne *regexp.Regexp


func init()  {
	var e error
	regUriAll, e = regexp.Compile(`/api/([a-zA-Z]+)/?([a-zA-Z/]+)?$`)
	if e != nil {
		panic(e)
	}

	regUriOne, e = regexp.Compile(`/api/([a-zA-Z]+)/([0-9]+)/?([a-zA-Z/]+)?$`)
	if e != nil {
		panic(e)
	}

}

type Server struct {
	Config lib1ssarp.Configuration
}


func (s Server) Launch() {
	listen := s.Config.Server.Address()
	log.Println("listen: ", listen)

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc( s.index ))

	log.Fatal(http.ListenAndServe(listen, mux))
}


func (s Server) index(w http.ResponseWriter, r *http.Request) {
	log.Println("Index page", r.URL.Path, r.Method);

	//simply route
	switch r.Method {
	case METHOD_GET:
		if regUriOne.MatchString(r.URL.Path) {
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
func (s Server) all(w http.ResponseWriter, r *http.Request) {

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

	//TODO delegate to service model
	//TODO return JSON response
	log.Println("Model Name: ", model, ",  relation: ", relation)
	renderJson(w, fmt.Sprintf(`[{"Model": "%s", "Relation": "%s"}]`, model, relation));
}

func (s Server) one(w http.ResponseWriter, r *http.Request) {

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

	//TODO delegate to service model
	//TODO return JSON response
	log.Println("Model Name: ", model, ", pk: ", id, ", relation: ", relation)
	renderJson(w, fmt.Sprintf(`{"Model": "%s", "Pk": %s, "Relation": "%s"}`, model, id, relation));

}

func (s Server) create(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Create, Path: %q, Method: %s", r.URL.Path, r.Method)
}

func (s Server) update(w http.ResponseWriter, r *http.Request) {

}

func (s Server) delete(w http.ResponseWriter, r *http.Request) {

}


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