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
	"sync"
	"time"

	"github.com/axelspringer/swerve/src/db"
	"golang.org/x/crypto/acme/autocert"
)

// Manager wraps around autocert and injects a cache
type Manager struct {
	certCache   *persistentCertCache
	acmeManager *autocert.Manager
}

// persistentCertCache certificate cache
type persistentCertCache struct {
	autocert.Cache
	db              *db.DynamoDB
	pollTicker      *time.Ticker
	mapMutex        *sync.Mutex
	wildcardDomains []*db.Domain
	domainsMap      map[string]*db.Domain
}
