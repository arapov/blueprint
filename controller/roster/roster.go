package roster

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/arapov/pile2/lib/flight"
	"github.com/arapov/pile2/model/roster"

	"github.com/blue-jay-fork/core/router"
)

// Must obey http://jsonapi.org/format/
type response struct {
	Data   []map[string][]string `json:"data"`
	Errors []map[string]string   `json:"errors"`
	Meta   map[string]string     `json:"meta"`
}

// Load the routes.
func Load() {
	router.Get("/roster", Index)
	router.Get("/roster/", Index)
	router.Get("/roster/:group", Index)

	router.Get("/api/v1/roster/groups", apiGetGroups)
	router.Get("/api/v1/roster/groups/", apiGetGroups)
	router.Get("/api/v1/roster/groups/:group", apiGetGroups)
}

func apiGetGroups(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	c := flight.Context(w, r)
	start := time.Now()

	var filter string
	var res response
	var data []map[string][]string

	data, err := roster.GetGroups(c.LDAP, c.Param("group"), &filter)
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
	}

	w.Write(jsonRes)
}

// Index displays the page.
func Index(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)
	v := c.View.New("roster/index")

	v.Vars["group"] = c.Param("group")

	v.Render(w, r)
}
