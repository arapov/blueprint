// Package flight_test
package flight_test

import (
	"log"
	"testing"

	"github.com/arapov/pile2/lib/flight"
	"github.com/blue-jay-fork/core/xsrf"
)

// TestRace tests for race conditions.
func TestXsrfRace(t *testing.T) {
	for i := 0; i < 100; i++ {
		go func() {
			// Set the csrf information
			flight.StoreXSRF(xsrf.Info{
				AuthKey: "test123",
				Secure:  true,
			})
			x := flight.XSRF()
			x.AuthKey = "test"
			log.Println(flight.XSRF().AuthKey)
		}()
	}
}
