package lineacaffe

import (
	"github.com/mikepartelow/coffeemenu"
)

type LineaCaffe struct {
	*coffeemenu.CollyScraper
}

func New() *LineaCaffe {
	return &LineaCaffe{
		CollyScraper: coffeemenu.NewCollyScraper(
			"Linea Caffe",
			[]string{
				"https://lineacaffe.com/shop/coffee/",
			},
			coffeemenu.ScrapeSpec{
				Container: "div.product",
				Name:      []string{"ChildText", "h4"},
				Url:       []string{"ChildAttr", "a", "href"},
			},
		),
	}
}
