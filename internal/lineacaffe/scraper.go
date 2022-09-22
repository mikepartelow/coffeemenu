package lineacaffe

import (
	"github.com/gocolly/colly/v2"
	"github.com/mikepartelow/coffeemenu"
)

type LineaCaffe struct {
	coffeemenu.CollyScraper
}

func New() *LineaCaffe {
	m := LineaCaffe{
		CollyScraper: coffeemenu.CollyScraper{
			Colly: colly.NewCollector(),
			Name:  "Linea Caffe",
			Urls: []string{
				"https://lineacaffe.com/shop/coffee/",
			},
		},
	}
	m.CollyScraper.Colly.OnHTML("div.product", func(e *colly.HTMLElement) {
		m.CollyScraper.Products = append(m.CollyScraper.Products, coffeemenu.Product{
			Name: e.ChildText("h4"),
			Url:  e.ChildAttr("a", "href"),
		})
	})

	return &m
}
