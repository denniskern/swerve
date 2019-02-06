package db_test

import (
	"errors"
	"net/url"
	"testing"

	"github.com/TetsuyaXD/swerve/src/db"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDBDomain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DB domain suite")
}

var _ = Describe("type Domain", func() {
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

		It("Domain struct redirect with query string not promotable and trailing slash", func() {
			domain := &db.Domain{
				Name:         "example.com",
				Redirect:     "https://www.example.com/",
				RedirectCode: 301,
				Promotable:   false,
			}

			url, err := url.Parse("https://example.com/path/to/promote?with=query")
			Expect(err).To(BeNil())

			redirectURL, redirectCode := domain.GetRedirect(url)
			Expect("https://www.example.com/").To(Equal(redirectURL))
			Expect(redirectCode).To(Equal(301))
		})

		It("Domain struct redirect without query string and not promotable", func() {
			domain := &db.Domain{
				Name:         "example.com",
				Redirect:     "https://www.example.com/path/to/promote/",
				RedirectCode: 301,
				Promotable:   false,
			}

			url, err := url.Parse("https://example.com/path/to/promote/")
			Expect(err).To(BeNil())

			redirectURL, redirectCode := domain.GetRedirect(url)
			Expect("https://www.example.com/path/to/promote/").To(Equal(redirectURL))
			Expect(redirectCode).To(Equal(301))
		})

		It("Domain struct redirect promotable without query string", func() {
			domain := &db.Domain{
				Name:         "example.com",
				Redirect:     "https://www.example.com",
				RedirectCode: 301,
				Promotable:   true,
			}

			url, err := url.Parse("https://example.com/path/to/promote/")
			Expect(err).To(BeNil())

			redirectURL, redirectCode := domain.GetRedirect(url)
			Expect("https://www.example.com/path/to/promote/").To(Equal(redirectURL))
			Expect(redirectCode).To(Equal(301))
		})

		It("Domain struct redirect promotable with query string", func() {
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
		It("Domain struct redirect path mapping and not promotable", func() {
			domain := &db.Domain{
				Name:         "example.com",
				Redirect:     "https://www.example.com/",
				RedirectCode: 301,
				Promotable:   false,
				PathMapping: &db.PathList{
					db.PathMappingEntry{From: "/old/path", To: "/new/target"},
				},
			}

			url, err := url.Parse("https://example.com/path/to/promote?with=query")
			Expect(err).To(BeNil())
			redirectURL, redirectCode := domain.GetRedirect(url)
			Expect("https://www.example.com/").To(Equal(redirectURL))
			Expect(redirectCode).To(Equal(301))

			url, err = url.Parse("https://example.com/old/path/to/promote?with=query")
			Expect(err).To(BeNil())
			redirectURL, redirectCode = domain.GetRedirect(url)
			Expect("https://www.example.com/new/target").To(Equal(redirectURL))
			Expect(redirectCode).To(Equal(301))

			mapping := *domain.PathMapping
			mapping[0].To = "https://theotherserver.com/new/target/"
			url, err = url.Parse("https://example.com/old/path/to/promote?with=query")
			Expect(err).To(BeNil())
			redirectURL, redirectCode = domain.GetRedirect(url)
			Expect("https://theotherserver.com/new/target/to/promote").To(Equal(redirectURL))
			Expect(redirectCode).To(Equal(301))
		})
	})
})
