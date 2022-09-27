package coffeemenu_test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/mikepartelow/coffeemenu"
	"github.com/stretchr/testify/assert"
)

type RoundTripFunc func(r *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestScraper(t *testing.T) {
	ProductName := "Colombia Supremo"
	ProductUrl := "/colombia-supremo/4"
	TestPage := fmt.Sprintf("<html><body><div><div class=\"name\">%s</div><div class=\"url\">%s</div></div></body></html>", ProductName, ProductUrl)

	rt := RoundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": {"text/html; charset=utf-8"}},
			Body:       io.NopCloser(strings.NewReader(TestPage)),
		}, nil
	})

	t.Run("it scrapes", func(t *testing.T) {
		site := coffeemenu.Site{
			Name: "GUE Coffee Roasters",
			Urls: []string{"https://gueroasters.com/shop"},
			ScrapeSpec: coffeemenu.ScrapeSpec{
				Container: "div",
				Name:      []string{"ChildText", ".name"},
				Url:       []string{"ChildText", ".url"},
			},
		}

		s := coffeemenu.NewScraper(site, rt)
		err := s.Scrape()
		assert.NoError(t, err)
		assert.Len(t, s.Products(), 1)

		p := s.Products()[0]
		assert.Equal(t, p.Name, ProductName)
		assert.Equal(t, p.Url, ProductUrl)
	})
}
