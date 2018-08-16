package gitpage

import (
	"net/http"

	"github.com/arapov/pile2/lib/flight"
	"github.com/arapov/pile2/model/gitpage"
	"github.com/blue-jay-fork/core/router"
)

// Load the routes.
func Load() {
	router.Get("/wiki/:page", Index)
}

// Index displays the home page.
func Index(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	v := c.View.New("wiki/index")
	v.Vars["title"] = c.Param("page")
	v.Vars["page"] = gitpage.Get(c.Git, c.Param("page"))
	v.Render(w, r)
}
