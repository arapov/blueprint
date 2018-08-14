package ldapxrest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/blue-jay-fork/blueprint/lib/flight"
	"github.com/blue-jay-fork/blueprint/model/ldapxrest"
	"github.com/blue-jay-fork/core/router"
)

var (
	uri = "/api"
)

// Must obey http://jsonapi.org/format/
type response struct {
	Data   []map[string][]string `json:"data"`
	Errors []map[string]string   `json:"errors"`
}

// Load the routes.
func Load() {
	router.Get(uri+"/v1/query", Query)
	router.Get(uri+"/v1/ping", Ping)
}

func Query(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var res response
	var data []map[string][]string

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		// TODO: use response struct
		w.Write([]byte(fmt.Sprintf("{\"response\":\"%s\"}", err)))

		return
	}

	c := flight.Context(w, r)
	data, err = ldapxrest.Query(c.LDAP, r.Form)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadGateway)
		// TODO: use response struct
		w.Write([]byte(fmt.Sprintf("{\"response\":\"%s\"}", err)))

		return
	}
	res.Data = data

	jsonRes, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		// TODO: use response struct
		w.Write([]byte(fmt.Sprintf("{\"response\":\"%s\"}", err)))

		return
	}

	w.Write(jsonRes)
}

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
