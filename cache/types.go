package cache

import (
	"sync"
	"time"

	"github.com/axelspringer/swerve/database"
)

// Cache contains the local redirect entry cache
type Cache struct {
	DB           DatabaseAdapter
	Observing    bool
	closer       chan struct{}
	mapMutex     *sync.RWMutex
	redirectsMap map[string]*database.Redirect
	certOrderMap map[string]time.Time
}
