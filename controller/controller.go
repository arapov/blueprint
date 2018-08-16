// Package controller loads the routes for each of the controllers.
package controller

import (
	"github.com/blue-jay-fork/blueprint/controller/about"
	"github.com/blue-jay-fork/blueprint/controller/debug"
	"github.com/blue-jay-fork/blueprint/controller/gitpage"
	"github.com/blue-jay-fork/blueprint/controller/home"
	"github.com/blue-jay-fork/blueprint/controller/ldapxrest"
	"github.com/blue-jay-fork/blueprint/controller/login"
	"github.com/blue-jay-fork/blueprint/controller/notepad"
	"github.com/blue-jay-fork/blueprint/controller/register"
	"github.com/blue-jay-fork/blueprint/controller/static"
	"github.com/blue-jay-fork/blueprint/controller/status"
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
}
