// Package env reads the application settings.
package env

import (
	"encoding/json"
	"os"
	"regexp"

	"github.com/blue-jay-fork/core/asset"
	"github.com/blue-jay-fork/core/email"
	"github.com/blue-jay-fork/core/form"
	"github.com/blue-jay-fork/core/generate"
	"github.com/blue-jay-fork/core/jsonconfig"
	"github.com/blue-jay-fork/core/server"
	"github.com/blue-jay-fork/core/session"
	"github.com/blue-jay-fork/core/storage/driver/mysql"
	"github.com/blue-jay-fork/core/view"
)

// *****************************************************************************
// Application Settings
// *****************************************************************************

// Info structures the application settings.
type Info struct {
	Asset      asset.Info    `json:"Asset"`
	Email      email.Info    `json:"Email"`
	Form       form.Info     `json:"Form"`
	Generation generate.Info `json:"Generation"`
	MySQL      mysql.Info    `json:"MySQL"`
	Server     server.Info   `json:"Server"`
	Session    session.Info  `json:"Session"`
	Template   view.Template `json:"Template"`
	View       view.Info     `json:"View"`
	path       string
}

// Path returns the env.json path
func (c *Info) Path() string {
	return c.path
}

// ParseJSON unmarshals bytes to structs
func (c *Info) ParseJSON(b []byte) error {

	// w is always "env:*" here, hence [4:]
	Getenv := func(w string) string {
		return os.Getenv(w[4:])
	}

	// Looking for all the "env:*" strings to replace with ENV vars
	r := regexp.MustCompile(`env:[A-Z_]+`)
	newB := r.ReplaceAllStringFunc(string(b), Getenv)

	return json.Unmarshal([]byte(newB), &c)
}

// New returns a instance of the application settings.
func New(path string) *Info {
	return &Info{
		path: path,
	}
}

// LoadConfig reads the configuration file.
func LoadConfig(configFile string) (*Info, error) {
	// Create a new configuration with the path to the file
	config := New(configFile)

	// Load the configuration file
	err := jsonconfig.Load(configFile, config)

	// Return the configuration
	return config, err
}
