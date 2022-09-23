package rebootroasting

import (
	"github.com/mikepartelow/coffeemenu"
)

type RebootRoasting struct {
	*coffeemenu.CollyScraper
}

func New() *RebootRoasting {
	return &RebootRoasting{
		CollyScraper: coffeemenu.NewCollyScraper(
			"Reboot Roasting",
			[]string{
				"https://www.rebootroasting.com/products?category=Coffee",
			},
			coffeemenu.ScrapeSpec{
				Container: "a.product",
				Name:      []string{"ChildText", "div.product-title"},
				Url:       []string{"Attr", "href"},
			},
		),
	}
}
