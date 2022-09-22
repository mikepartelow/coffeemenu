package rebootroasting

import (
	"github.com/gocolly/colly/v2"
	"github.com/mikepartelow/coffeemenu"
)

type RebootRoasting struct {
	coffeemenu.CollyScraper
}

func New() *RebootRoasting {
	m := RebootRoasting{
		CollyScraper: coffeemenu.CollyScraper{
			Colly: colly.NewCollector(),
			Name:  "Reboot Roasting",
			Urls: []string{
				"https://www.rebootroasting.com/products?category=Coffee",
			},
		},
	}
	m.CollyScraper.Colly.OnHTML("a.product", func(e *colly.HTMLElement) {
		m.CollyScraper.Products = append(m.CollyScraper.Products, coffeemenu.Product{
			Name: e.ChildText("div.product-title"),
			Url:  e.Attr("href"),
		})
	})

	return &m
}
