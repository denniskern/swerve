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

package configuration_test

import (
	"errors"
	"testing"

	"github.com/TetsuyaXD/swerve/src/db"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConfiguration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Configuration Suite")
}

var _ = Describe("Configuration", func() {
	It("Domain struct validating", func() {
		domain := &db.Domain{}
		errList := domain.Validate()
		Expect(errList).To(Equal([]error{
			errors.New("Invalid id"),
			errors.New("Invalid domain name"),
			errors.New("Invalid domain date"),
			errors.New("Invalid domain redirect target"),
			errors.New("Invalid redirect http status code"),
		}))
	})
})
