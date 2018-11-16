package configuration_test

import (
	"errors"
	"testing"

	"github.com/axelspringer/swerve/src/db"

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
