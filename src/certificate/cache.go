// Copyright 2018 Axel Springer SE
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package certificate

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/axelspringer/swerve/src/db"
	"github.com/axelspringer/swerve/src/log"
	"golang.org/x/crypto/acme/autocert"
)

const (
	pollTickerInterval = 10
)

// newPersistentCertCache creates a new persistent cache based on dynamo db
func newPersistentCertCache(d *db.DynamoDB) *persistentCertCache {
	c := &persistentCertCache{
		pollTicker:      time.NewTicker(time.Second * pollTickerInterval),
		db:              d,
		domainsMap:      map[string]*db.Domain{},
		wildcardDomains: []*db.Domain{},
		mapMutex:        &sync.Mutex{},
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

	log.Debugf("persistentCertCache.updateDomainCache %#v", domains)

	m := map[string]*db.Domain{}
	w := []*db.Domain{}
	for _, domain := range domains {
		m[domain.Name] = &domain
		if domain.Wildcard == true {
			w = append(w, &domain)
		}
	}

	c.mapMutex.Lock()
	c.domainsMap = m
	c.wildcardDomains = w
	c.mapMutex.Unlock()
}

// Get cert by domain name
func (c *persistentCertCache) Get(ctx context.Context, key string) ([]byte, error) {
	c.mapMutex.Lock()
	defer c.mapMutex.Unlock()

	// check for non wildcard domains
	if domain, ok := c.domainsMap[key]; ok {
		if len(domain.Certificate) > 0 {
			return []byte(domain.Certificate), nil
		}
	}

	// check wildcard domains
	for _, wc := range c.wildcardDomains {
		if strings.HasSuffix(key, "."+wc.Name) {
			if len(wc.Certificate) > 0 {
				return []byte(wc.Certificate), nil
			}
			return nil, autocert.ErrCacheMiss
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
		log.Debugf("persistentCertCache.Put %s %#v", key, data)
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
		log.Debugf("persistentCertCache.Delete %s ''", key)
		err = c.db.UpdateCertificateData(key, []byte{})
		close(done)
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
			log.Debug("Update domain cache")
			c.updateDomainCache()
		}
	}()

	return nil
}

// IsDomainAcceptable test for domains in cache
func (c *persistentCertCache) IsDomainAcceptable(domain string) (*db.Domain, bool) {
	// check non wildcard domains
	c.mapMutex.Lock()
	if d, ok := c.domainsMap[domain]; ok {
		return d, ok
	}
	c.mapMutex.Unlock()

	// check wildcard domains
	for _, wc := range c.wildcardDomains {
		if strings.HasSuffix(domain, "."+wc.Name) {
			return wc, true
		}
	}

	return nil, false
}
