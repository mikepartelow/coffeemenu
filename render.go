package coffeemenu

import (
	"bytes"
	"encoding/csv"
	"io"
	"sort"
	"text/template"

	_ "embed"

	"github.com/charmbracelet/glamour"
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
		if err := tmpl.Execute(w, s); err != nil {
			log.Error().Err(err).Send()
		}
	}
}

func Glamourize(in string) string {
	r, err := glamour.NewTermRenderer(
		glamour.WithStylePath("dracula"),
		glamour.WithWordWrap(0),
	)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	out, err := r.Render(in)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	return out
}

func CSV(scrapers []*Scraper) string {
	var buf bytes.Buffer

	w := csv.NewWriter(&buf)
	_ = w.Write([]string{"product", "roaster"})

	for _, s := range scrapers {
		log.Debug().Msg(s.Name())
		for _, p := range s.Products() {
			_ = w.Write([]string{p.Name, s.Name()})
		}
	}
	w.Flush()
	return buf.String()
}
