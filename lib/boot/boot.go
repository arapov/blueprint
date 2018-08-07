// Package boot handles the initialization of the web components.
package boot

import (
	"log"

	"github.com/blue-jay-fork/blueprint/controller"
	"github.com/blue-jay-fork/blueprint/lib/env"
	"github.com/blue-jay-fork/blueprint/lib/flight"
	"github.com/blue-jay-fork/blueprint/viewfunc/link"
	"github.com/blue-jay-fork/blueprint/viewfunc/noescape"
	"github.com/blue-jay-fork/blueprint/viewfunc/prettytime"
	"github.com/blue-jay-fork/blueprint/viewmodify/authlevel"
	"github.com/blue-jay-fork/blueprint/viewmodify/flash"
	"github.com/blue-jay-fork/blueprint/viewmodify/uri"

	"github.com/blue-jay-fork/core/form"
	"github.com/blue-jay-fork/core/pagination"
	"github.com/blue-jay-fork/core/xsrf"
)

// RegisterServices sets up all the web components.
func RegisterServices(config *env.Info) {
	// Set up the session cookie store
	err := config.Session.SetupConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Connect to the MySQL database
	mysqlDB, _ := config.MySQL.Connect(true)

	// Load the controller routes
	controller.LoadRoutes()

	// Set up the views
	config.View.SetTemplates(config.Template.Root, config.Template.Children)

	// Set up the functions for the views
	config.View.SetFuncMaps(
		config.Asset.Map(config.View.BaseURI),
		link.Map(config.View.BaseURI),
		noescape.Map(),
		prettytime.Map(),
		form.Map(),
		pagination.Map(),
	)

	// Set up the variables and modifiers for the views
	config.View.SetModifiers(
		authlevel.Modify,
		uri.Modify,
		xsrf.Token,
		flash.Modify,
	)

	// Store the variables in flight
	flight.StoreConfig(*config)

	// Store the database connection in flight
	flight.StoreDB(mysqlDB)

	// Store the csrf information
	flight.StoreXSRF(xsrf.Info{
		AuthKey: config.Session.CSRFKey,
		Secure:  config.Session.Options.Secure,
	})
}
