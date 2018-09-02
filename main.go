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

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/axelspringer/swerve/src/app"
)

// Version filled by the build process
var Version string

func main() {
	// new application
	application := app.NewApplication()
	// read configuration
	application.Setup()
	// print app verion
	if application.Config.Version {
		fmt.Printf("version %v\n", Version)
		os.Exit(0)
	}
	// print app config
	if application.Config.Help {
		flag.PrintDefaults()
		os.Exit(0)
	}
	// run the server
	application.Run()
}
