package coffeemenu_test

import (
	"bytes"
	"embed"
	_ "embed"
	"io"
	"net/http"
	"testing"

	"github.com/mikepartelow/coffeemenu"
	"github.com/stretchr/testify/assert"
)

type RoundTripFunc func(r *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

//go:embed fixtures
var fixtures embed.FS

func TestScraper(t *testing.T) {
	testCases := []struct {
		productName   string
		productUrl    string
		fixture       string
		containerSpec string
		nameSpec      []string
		urlSpec       []string
	}{
		{
			productName:   "Colombia Supremo",
			productUrl:    "/colombia-supremo/4",
			fixture:       "scraper.0.html",
			containerSpec: "div",
			nameSpec:      []string{"ChildText", ".name"},
			urlSpec:       []string{"ChildText", ".url"},
		},
		{
			productName:   "Ethiopia Sidamo",
			productUrl:    "/ethiopia-sidamo/9",
			fixture:       "scraper.1.html",
			containerSpec: "div",
			nameSpec:      []string{"ChildAttr", ".name", "bean"},
			urlSpec:       []string{"Attr", "url"},
		},
		{
			productName:   "Mexico Chiapas",
			productUrl:    "mexico-chiapas",
			fixture:       "scraper.2.html",
			containerSpec: "div",
			nameSpec:      []string{"ChildText", ".name"},
			urlSpec:       []string{"ChildAttrBase", "a", "href"},
		},
	}

	t.Run("it scrapes", func(t *testing.T) {
		for _, tC := range testCases {
			t.Run(tC.fixture, func(t *testing.T) {
				testPage, err := fixtures.ReadFile("fixtures/" + tC.fixture)
				assert.NoError(t, err)

				rt := makeRoundTripper(testPage)
				site := makeSite(tC.containerSpec, tC.nameSpec, tC.urlSpec)

				s := coffeemenu.NewScraper(site, rt)
				err = s.Scrape()
				assert.NoError(t, err)
				assert.Len(t, s.Products(), 1)

				p := s.Products()[0]
				assert.Equal(t, p.Name, tC.productName)
				assert.Equal(t, p.Url, tC.productUrl)
			})
		}
	})
}

func makeRoundTripper(testPage []byte) RoundTripFunc {
	return RoundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": {"text/html; charset=utf-8"}},
			Body:       io.NopCloser(bytes.NewReader(testPage)),
		}, nil
	})
}

func makeSite(container string, name, url []string) coffeemenu.Site {
	return coffeemenu.Site{
		Name: "GUE Coffee Roasters",
		Urls: []string{"https://gueroasters.com/shop"},
		ScrapeSpec: coffeemenu.ScrapeSpec{
			Container: container,
			Name:      name,
			Url:       url,
		},
	}

}
