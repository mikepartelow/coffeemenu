package coffeemenu

import (
	"github.com/gocolly/colly/v2"
)

type Scraper interface {
	Scrape() error
	GetName() string
	GetProducts() Products
	GetStats() string
}

type CollyScraper struct {
	Colly    *colly.Collector
	Name     string
	Products Products
	Urls     []string
}

func (cs CollyScraper) GetName() string {
	return cs.Name
}

func (cs CollyScraper) GetProducts() Products {
	return cs.Products
}

func (cs CollyScraper) Scrape() error {
	for _, url := range cs.Urls {
		if err := cs.Colly.Visit(url); err != nil {
			return err
		}
	}

	return nil
}

func (cs CollyScraper) GetStats() string {
	return cs.Colly.String()
}
