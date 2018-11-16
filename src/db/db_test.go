package db_test

import (
	"errors"
	"net/url"
	"testing"

	"github.com/axelspringer/swerve/src/db"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDB(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DB Suite")
}

var _ = Describe("DB", func() {
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

	Context("type Domain struct", func() {
		It("Domain struct redirect", func() {
			domain := &db.Domain{
				Name:         "example.com",
				Redirect:     "https://www.example.com",
				RedirectCode: 301,
			}

			url, err := url.Parse("https://example.com/path/to/promote?with=query")
			Expect(err).To(BeNil())

			redirectURL, redirectCode := domain.GetRedirect(url)
			Expect("https://www.example.com").To(Equal(redirectURL))
			Expect(redirectCode).To(Equal(301))
		})

		It("Domain struct redirect promotable", func() {
			domain := &db.Domain{
				Name:         "example.com",
				Redirect:     "https://www.example.com",
				RedirectCode: 301,
				Promotable:   true,
			}

			url, err := url.Parse("https://example.com/path/to/promote?with=query")
			Expect(err).To(BeNil())

			redirectURL, redirectCode := domain.GetRedirect(url)
			Expect("https://www.example.com/path/to/promote?with=query").To(Equal(redirectURL))
			Expect(redirectCode).To(Equal(301))
		})

		It("Domain struct redirect promotable path mapping", func() {
			domain := &db.Domain{
				Name:         "example.com",
				Redirect:     "https://www.example.com",
				RedirectCode: 301,
				Promotable:   true,
				PathMapping: &db.PathList{
					db.PathMappingEntry{From: "/old/path", To: "/new/target"},
				},
			}

			url, err := url.Parse("https://example.com/path/to/promote?with=query")
			Expect(err).To(BeNil())
			redirectURL, redirectCode := domain.GetRedirect(url)
			Expect("https://www.example.com/path/to/promote?with=query").To(Equal(redirectURL))
			Expect(redirectCode).To(Equal(301))

			url, err = url.Parse("https://example.com/old/path/to/promote?with=query")
			Expect(err).To(BeNil())
			redirectURL, redirectCode = domain.GetRedirect(url)
			Expect("https://www.example.com/new/target/to/promote?with=query").To(Equal(redirectURL))
			Expect(redirectCode).To(Equal(301))

			mapping := *domain.PathMapping
			mapping[0].To = "https://theotherserver.com/new/target/"
			url, err = url.Parse("https://example.com/old/path/to/promote?with=query")
			Expect(err).To(BeNil())
			redirectURL, redirectCode = domain.GetRedirect(url)
			Expect("https://theotherserver.com/new/target/to/promote?with=query").To(Equal(redirectURL))
			Expect(redirectCode).To(Equal(301))
		})
	})
})
