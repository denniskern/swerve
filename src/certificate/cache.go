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
	"sync"
	"time"

	"github.com/TetsuyaXD/swerve/src/db"
	"github.com/TetsuyaXD/swerve/src/log"
	"golang.org/x/crypto/acme/autocert"
)

const (
	pollTickerInterval = 10
)

// NewPersistentCertCache creates a new persistent cache based on dynamo db
func NewPersistentCertCache(d *db.DynamoDB) *PersistentCertCache {
	return &PersistentCertCache{
		PollTicker: time.NewTicker(time.Second * pollTickerInterval),
		DB:         d,
		DomainsMap: map[string]db.Domain{},
		MapMutex:   &sync.Mutex{},
	}
}

// UpdateDomainCache updates the domain cache
func (c *PersistentCertCache) UpdateDomainCache() {
	domains, err := c.DB.FetchAll()
	if err != nil {
		log.Errorf("Error while fetching domain list %v", err)
		return
	}

	// lock the map
	c.MapMutex.Lock()
	defer c.MapMutex.Unlock()
	// create new domain map
	c.DomainsMap = map[string]db.Domain{}

	for _, domain := range domains {
		c.DomainsMap[domain.Name] = domain
	}
}

// Get cert by domain name
func (c *PersistentCertCache) Get(ctx context.Context, key string) ([]byte, error) {
	var (
		done = make(chan struct{})
		err  error
		data []byte
	)

	go func() {
		defer close(done)
		data, err = c.DB.GetTLSCache(key)
	}()

	// handle context timeouts and errors
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-done:
	}

	if err == nil && data != nil && len(data) > 0 {
		return data, nil
	}

	return nil, autocert.ErrCacheMiss
}

// Put a cert to the cache
func (c *PersistentCertCache) Put(ctx context.Context, key string, data []byte) error {
	var (
		done = make(chan struct{})
		err  error
	)

	go func() {
		defer close(done)
		err = c.DB.UpdateTLSCache(key, data)
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
func (c *PersistentCertCache) Delete(ctx context.Context, key string) error {
	var (
		done = make(chan struct{})
		err  error
	)

	go func() {
		defer close(done)
		err = c.DB.DeleteTLSCacheEntry(key)
	}()

	// handle context timeouts and errors
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
	}

	return err
}

// Observe the domain backend. Ya through polling. Pub/Sub would be much better. Go implement it
func (c *PersistentCertCache) Observe() error {
	go func() {
		for _ = range c.PollTicker.C {
			c.UpdateDomainCache()
		}
	}()

	return nil
}

// IsDomainAcceptable test for domains in cache
func (c *PersistentCertCache) IsDomainAcceptable(domain string) (*db.Domain, bool) {
	// check non wildcard domains
	c.MapMutex.Lock()
	defer c.MapMutex.Unlock()

	if d, ok := c.DomainsMap[domain]; ok {
		return &d, ok
	}

	return nil, false
}
