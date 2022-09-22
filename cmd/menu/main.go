package main

import (
	"html/template"
	"io"
	"os"
	"sort"
	"sync"

	_ "embed"

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

	for _, s := range scrapers {
		if err := tmpl.Execute(os.Stdout, s); err != nil {
			log.Error().Err(err).Send()
		}
	}
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
		s := s
		go func() {
			if err := s.Scrape(); err != nil {
				log.Error().Err(err).Msgf("Error while scraping %q", s.GetName())
			}
			wg.Done()
		}()
	}

	wg.Wait()

	render(scrapers, os.Stdout)
}
