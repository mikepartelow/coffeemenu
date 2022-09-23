package coffeemenu

import (
	"io"
	"os"
	"sort"
	"text/template"

	_ "embed"

	"github.com/rs/zerolog/log"
)

//go:embed menu.md.tmpl
var menu string

func Render(scrapers []*Scraper, w io.Writer) {
	sorted := func(products Products) Products {
		sort.Sort(products)
		return products
	}

	t := template.New("menu")
	tmpl := template.Must(t.Funcs(template.FuncMap{"sorted": sorted}).Parse(menu))

	for _, s := range scrapers {
		if err := tmpl.Execute(os.Stdout, s); err != nil {
			log.Error().Err(err).Send()
		}
	}
}
