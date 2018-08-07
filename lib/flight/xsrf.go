package flight

import (
	"sync"

	"github.com/blue-jay-fork/core/xsrf"
)

var (
	xsrfInfo  xsrf.Info
	xsrfMutex sync.RWMutex
)

// StoreXSRF sets the csrf configuration.
func StoreXSRF(x xsrf.Info) {
	xsrfMutex.Lock()
	xsrfInfo = x
	xsrfMutex.Unlock()
}

// XSRF returns the csrf configuration.
func XSRF() xsrf.Info {
	xsrfMutex.RLock()
	x := xsrfInfo
	xsrfMutex.RUnlock()
	return x
}
