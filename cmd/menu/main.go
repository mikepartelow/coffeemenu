package main

import (
	"bytes"
	"io"
	"os"
	"sort"
	"sync"
	"text/template"

	_ "embed"

	"github.com/charmbracelet/glamour"
	"github.com/mikepartelow/coffeemenu"
	"github.com/mikepartelow/coffeemenu/internal/lineacaffe"
	"github.com/mikepartelow/coffeemenu/internal/mrespresso"
	"github.com/mikepartelow/coffeemenu/internal/rebootroasting"
	"github.com/rs/zerolog/log"
)

//go:embed menu.md.tmpl
var menu string

func render(scrapers []coffeemenu.Scraper, w io.Writer) {
	sorted := func(products coffeemenu.Products) coffeemenu.Products {
		sort.Sort(products)
		return products
	}

	t := template.New("menu")
	tmpl := template.Must(t.Funcs(template.FuncMap{"sorted": sorted}).Parse(menu))

	var buf bytes.Buffer

	for _, s := range scrapers {
		if err := tmpl.Execute(&buf, s); err != nil {
			log.Error().Err(err).Send()
		}
	}

	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(0),
	)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	out, err := r.Render(buf.String())
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	_, _ = w.Write([]byte(out))
}

func main() {
	scrapers := []coffeemenu.Scraper{
		mrespresso.New(),
		rebootroasting.New(),
		lineacaffe.New(),
	}

	var wg sync.WaitGroup

	for _, s := range scrapers {
		wg.Add(1)
		go func(s coffeemenu.Scraper) {
			if err := s.Scrape(); err != nil {
				log.Error().Err(err).Msgf("Error while scraping %q", s.GetName())
			}
			wg.Done()
		}(s)
	}

	wg.Wait()

	render(scrapers, os.Stdout)
}
