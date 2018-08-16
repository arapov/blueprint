// Package controller loads the routes for each of the controllers.
package controller

import (
	"github.com/arapov/pile2/controller/about"
	"github.com/arapov/pile2/controller/debug"
	"github.com/arapov/pile2/controller/gitpage"
	"github.com/arapov/pile2/controller/home"
	"github.com/arapov/pile2/controller/ldapxrest"
	"github.com/arapov/pile2/controller/login"
	"github.com/arapov/pile2/controller/notepad"
	"github.com/arapov/pile2/controller/register"
	"github.com/arapov/pile2/controller/roster"
	"github.com/arapov/pile2/controller/static"
	"github.com/arapov/pile2/controller/status"
)

// LoadRoutes loads the routes for each of the controllers.
func LoadRoutes() {
	about.Load()
	debug.Load()
	register.Load()
	login.Load()
	home.Load()
	static.Load()
	status.Load()
	notepad.Load()
	ldapxrest.Load()
	gitpage.Load()
	roster.Load()
}
