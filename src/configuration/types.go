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

package configuration

import "github.com/axelspringer/swerve/src/db"

// Version string
var Version string

// Configuration model
type Configuration struct {
	HTTPListener  string
	HTTPSListener string
	APIListener   string
	DynamoDB      db.DynamoConnection
	TablePrefix   string
	LogLevel      string
	LogFormatter  string
	Bootstrap     bool
	Version       bool
	Help          bool
	StagingCA     bool
	APISecret     string
}
