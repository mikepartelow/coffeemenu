package coffeemenu

import (
	"path/filepath"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"
)

type ProductSet map[string]Product

type Scraper struct {
	colly    *colly.Collector
	name     string
	products ProductSet
	urls     []string
}

func NewScraper(site Site) *Scraper {
	s := Scraper{
		colly:    colly.NewCollector(),
		name:     site.Name,
		urls:     site.Urls,
		products: make(ProductSet),
	}

	s.colly.OnHTML(site.ScrapeSpec.Container, func(e *colly.HTMLElement) {
		name := eFunc(e, site.ScrapeSpec.Name)
		s.products[name] = Product{
			Name: name,
			Url:  eFunc(e, site.ScrapeSpec.Url),
		}
	})

	return &s
}

func (s Scraper) Name() string {
	return s.name
}

func (s Scraper) Products() Products {
	var products Products
	for _, p := range s.products {
		products = append(products, p)
	}
	return products
}

func (s Scraper) Scrape() error {
	for _, url := range s.urls {
		if err := s.colly.Visit(url); err != nil {
			return err
		}
	}

	return nil
}

func (s Scraper) Stats() string {
	return s.colly.String()
}

func eFunc(e *colly.HTMLElement, args []string) string {
	switch args[0] {
	case "ChildText":
		return e.ChildText(args[1])
	case "ChildAttr":
		return e.ChildAttr(args[1], args[2])
	case "ChildAttrBase":
		return filepath.Base(e.ChildAttr(args[1], args[2]))
	case "Attr":
		return e.Attr(args[1])
	default:
		log.Fatal().Msgf("unknown nameFunc name: %q", args[0])
	}
	return ""
}
