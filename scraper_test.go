package coffeemenu_test

import (
	"bytes"
	"embed"
	_ "embed"
	"io"
	"net/http"
	"strings"
	"testing"
	"text/template"

	"github.com/mikepartelow/coffeemenu"
	"github.com/stretchr/testify/assert"
)

type RoundTripFunc func(r *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

type TestCase struct {
	ProductName   string
	ProductUrl    string
	Extra         string
	Join          string
	Fixture       string
	ContainerSpec string
	NameSpec      []string
	UrlSpec       []string
}

func TestScraper(t *testing.T) {
	testCases := []TestCase{
		{
			ProductName:   "Colombia Supremo",
			ProductUrl:    "/colombia-supremo/4",
			Fixture:       "scraper.0.html",
			ContainerSpec: "div",
			NameSpec:      []string{"ChildText", ".name"},
			UrlSpec:       []string{"ChildText", ".url"},
		},
		{
			ProductName:   "Ethiopia Sidamo",
			ProductUrl:    "/ethiopia-sidamo/9",
			Fixture:       "scraper.1.html",
			ContainerSpec: "div",
			NameSpec:      []string{"ChildAttr", ".name", "bean"},
			UrlSpec:       []string{"Attr", "url"},
		},
		{
			ProductName:   "Mexico Chiapas",
			ProductUrl:    "mexico-chiapas",
			Fixture:       "scraper.2.html",
			ContainerSpec: "div",
			NameSpec:      []string{"ChildText", ".name"},
			UrlSpec:       []string{"ChildAttrBase", "a", "href"},
		},
		{
			ProductName:   "Sulawesi Toraja",
			Extra:         "Mike's Estate",
			Join:          " - ",
			ProductUrl:    "sulawesi/toraja/show.html",
			Fixture:       "scraper.3.html",
			ContainerSpec: "div",
			NameSpec:      []string{"ChildTextJoin", ".name", ".extra", " - "},
			UrlSpec:       []string{"ChildAttr", "a", "href"},
		},
	}

	t.Run("it scrapes", func(t *testing.T) {
		for _, tC := range testCases {
			t.Run(tC.Fixture, func(t *testing.T) {
				testPage := makeTestPage(t, tC)

				rt := makeRoundTripper(testPage)
				site := makeSite(tC.ContainerSpec, tC.NameSpec, tC.UrlSpec)

				s := coffeemenu.NewScraper(site, rt)
				err := s.Scrape()
				assert.NoError(t, err)
				assert.Len(t, s.Products(), 1)

				wantProductName := tC.ProductName
				if tC.Extra != "" {
					wantProductName = strings.Join([]string{tC.ProductName, tC.Extra}, tC.Join)
				}

				p := s.Products()[0]
				assert.Equal(t, wantProductName, p.Name)
				assert.Equal(t, tC.ProductUrl, p.Url)
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

//go:embed fixtures
var fixtures embed.FS

func makeTestPage(t *testing.T, tC TestCase) []byte {
	testPage := &bytes.Buffer{}

	fixture, err := fixtures.ReadFile("fixtures/" + tC.Fixture)
	assert.NoError(t, err)

	err = template.Must(template.New(tC.Fixture).Parse(string(fixture))).Execute(testPage, tC)
	assert.NoError(t, err)

	return testPage.Bytes()
}
