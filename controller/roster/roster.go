package roster

import (
	"net/http"

	"github.com/arapov/pile2/lib/flight"

	"github.com/blue-jay-fork/core/router"
)

// Load the routes.
func Load() {
	router.Get("/roster", Index)
}

// Index displays the page.
func Index(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	v := c.View.New("roster/index")
	v.Render(w, r)
}
