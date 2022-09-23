package coffeemenu

import (
	"github.com/rs/zerolog/log"

	"github.com/gocolly/colly/v2"
)

type Scraper interface {
	Scrape() error
	Name() string
	Products() Products
	Stats() string
}

type CollyScraper struct {
	colly    *colly.Collector
	name     string
	products Products
	urls     []string
}

type ScrapeSpec struct {
	Container string
	Name      []string
	Url       []string
}

func NewCollyScraper(name string, urls []string, spec ScrapeSpec) *CollyScraper {
	s := CollyScraper{
		colly: colly.NewCollector(),
		name:  name,
		urls:  urls,
	}

	s.colly.OnHTML(spec.Container, func(e *colly.HTMLElement) {
		s.products = append(s.products, Product{
			Name: eFunc(e, spec.Name),
			Url:  eFunc(e, spec.Url),
		})
	})

	return &s
}

func (cs CollyScraper) Name() string {
	return cs.name
}

func (cs CollyScraper) Products() Products {
	return cs.products
}

func (cs CollyScraper) Scrape() error {
	for _, url := range cs.urls {
		if err := cs.colly.Visit(url); err != nil {
			return err
		}
	}

	return nil
}

func (cs CollyScraper) Stats() string {
	return cs.colly.String()
}

func eFunc(e *colly.HTMLElement, args []string) string {
	switch args[0] {
	case "ChildText":
		return e.ChildText(args[1])
	case "ChildAttr":
		return e.ChildAttr(args[1], args[2])
	case "Attr":
		return e.Attr(args[1])
	default:
		log.Fatal().Msgf("unknown nameFunc name: %q", args[0])
	}
	return ""
}
