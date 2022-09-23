package mrespresso

import (
	"github.com/mikepartelow/coffeemenu"
)

type MrEspresso struct {
	*coffeemenu.CollyScraper
}

func New() *MrEspresso {
	return &MrEspresso{
		CollyScraper: coffeemenu.NewCollyScraper(
			"Mr. Espresso",
			[]string{
				"https://mrespresso.com/shop/coffee/espresso/",
				"https://mrespresso.com/shop/coffee/single-origin/",
			},
			coffeemenu.ScrapeSpec{
				Container: "h3.product-title",
				Name:      []string{"ChildText", "a"},
				Url:       []string{"ChildAttr", "a", "href"},
			},
		),
	}
}
