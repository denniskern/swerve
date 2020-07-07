package cache

import (
	"sync"

	"github.com/TetsuyaXD/evade/database"
)

// Cache contains the local redirect entry cache
type Cache struct {
	DB           DatabaseAdapter
	Observing    bool
	closer       chan struct{}
	mapMutex     *sync.RWMutex
	redirectsMap map[string]*database.Redirect
}
