package ldapxrest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/arapov/pile2/lib/flight"
	"github.com/arapov/pile2/model/ldapxrest"
	"github.com/blue-jay-fork/core/router"
)

var (
	uri = "/api"
)

// Must obey http://jsonapi.org/format/
type response struct {
	Data   []map[string][]string `json:"data"`
	Errors []map[string]string   `json:"errors"`
	Meta   map[string]string     `json:"meta"`
}

// Load the routes.
func Load() {
	router.Get(uri+"/v1/query", Query)
	router.Get(uri+"/v1/ping", Ping)
}

// Query is general purpose LDAP search request
func Query(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	c := flight.Context(w, r)
	start := time.Now()

	var filter string
	var res response
	var data []map[string][]string

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		res.Errors = append(res.Errors, map[string]string{"title": err.Error()})

		goto out
	}

	data, err = ldapxrest.Query(c.LDAP, r.Form, &filter)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadGateway)
		res.Errors = append(res.Errors, map[string]string{"title": err.Error()})

		goto out
	}
	res.Data = data

out:
	elapsed := time.Since(start)
	res.Meta = map[string]string{
		"time":   elapsed.String(),
		"filter": filter,
	}
	jsonRes, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		res.Errors = append(res.Errors, map[string]string{"title": err.Error()})

		goto out
	}

	w.Write(jsonRes)
}

// Ping ensures we are connected to LDAP and able to query
func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/vnd.api+json")

	c := flight.Context(w, r)
	err := ldapxrest.Ping(c.LDAP)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadGateway)
		// TODO: use response struct
		w.Write([]byte(fmt.Sprintf("{\"response\":\"%s\"}", err)))
	} else {
		// TODO: use response struct
		w.Write([]byte("{\"response\":\"pong\"}"))
	}
}
