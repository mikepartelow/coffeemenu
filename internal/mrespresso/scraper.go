package mrespresso

import (
	"github.com/gocolly/colly/v2"
	"github.com/mikepartelow/coffeemenu"
)

type MrEspresso struct {
	coffeemenu.CollyScraper
}

func New() *MrEspresso {
	m := MrEspresso{
		CollyScraper: coffeemenu.CollyScraper{
			Colly: colly.NewCollector(),
			Name:  "Mr. Espresso",
			Urls: []string{
				"https://mrespresso.com/shop/coffee/espresso/",
				"https://mrespresso.com/shop/coffee/single-origin/",
			},
		},
	}
	m.CollyScraper.Colly.OnHTML("h3.product-title", func(e *colly.HTMLElement) {
		m.CollyScraper.Products = append(m.CollyScraper.Products, coffeemenu.Product{
			Name: e.ChildText("a"),
			Url:  e.ChildAttr("a", "href"),
		})
	})

	return &m
}
