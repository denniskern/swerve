package certificate

import (
	"context"
	"sync"
	"time"

	"github.com/hammi85/swerve/src/db"
	"github.com/hammi85/swerve/src/log"
	"golang.org/x/crypto/acme/autocert"
)

const (
	pollTickerInterval = 10
)

// newPersistentCertCache creates a new persistent cache based on dynamo db
func newPersistentCertCache(d *db.DynamoDB) *persistentCertCache {
	c := &persistentCertCache{
		pollTicker: time.NewTicker(time.Second * pollTickerInterval),
		db:         d,
		domainsMap: map[string]db.Domain{},
		mapMutex:   &sync.Mutex{},
	}

	// cache preload
	c.updateDomainCache()
	// backgroud update ticker
	c.observe()
	return c
}

// updateDomainCache updates the domain cache
func (c *persistentCertCache) updateDomainCache() {
	domains, err := c.db.FetchAll()
	if err != nil {
		log.Errorf("Error while fetching domain list %v", err)
		return
	}

	m := map[string]db.Domain{}
	for _, domain := range domains {
		m[domain.Name] = domain
	}

	c.mapMutex.Lock()
	c.domainsMap = m
	c.mapMutex.Unlock()
}

// Get cert by domain name
func (c *persistentCertCache) Get(ctx context.Context, key string) ([]byte, error) {
	c.mapMutex.Lock()
	defer c.mapMutex.Unlock()

	if domain, ok := c.domainsMap[key]; ok {
		if len(domain.Certificate) > 0 {
			return []byte(domain.Certificate), nil
		}
	}

	return nil, autocert.ErrCacheMiss
}

// Put a cert to the cache
func (c *persistentCertCache) Put(ctx context.Context, key string, data []byte) error {
	var (
		done = make(chan struct{})
		err  error
	)

	go func() {
		defer close(done)
		err = c.db.UpdateCertificateData(key, data)
	}()

	// handle context timeouts and errors
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
	}

	return err
}

// Delete a domain from
func (c *persistentCertCache) Delete(ctx context.Context, key string) error {
	var (
		done = make(chan struct{})
		err  error
	)

	go func() {
		defer close(done)
		err = c.db.UpdateCertificateData(key, []byte{})
	}()

	// handle context timeouts and errors
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
	}

	return err
}

// observe the domain backend. Ya through polling. Pub/Sub would be much better. Go implement it
func (c *persistentCertCache) observe() error {
	go func() {
		for _ = range c.pollTicker.C {
			c.updateDomainCache()
		}
	}()

	return nil
}

// IsDomainAcceptable test for domains in cache
func (c *persistentCertCache) IsDomainAcceptable(domain string) (*db.Domain, bool) {
	c.mapMutex.Lock()
	defer c.mapMutex.Unlock()

	d, ok := c.domainsMap[domain]
	return &d, ok
}
