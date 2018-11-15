package certificate_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCache(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cache Suite")
}

var _ = Describe("Cache", func() {
	Context("Domain cache", func() {

		It("Domain Cache lookup", func() {

		})

	})
})
