package main

import (
	"io"
	"os"
	"sort"
	"sync"
	"text/template"

	_ "embed"

	"github.com/mikepartelow/coffeemenu"
	"github.com/rs/zerolog/log"
)

//go:embed menu.md.tmpl
var menu string

func render(scrapers []*coffeemenu.Scraper, w io.Writer) {
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
	scrapers, err := coffeemenu.ReadScrapers()
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't read scrapers.")
	}

	var wg sync.WaitGroup

	for _, s := range scrapers {
		wg.Add(1)
		go func(s *coffeemenu.Scraper) {
			if err := s.Scrape(); err != nil {
				log.Error().Err(err).Msgf("Error while scraping %q", s.Name())
			}
			wg.Done()
		}(s)
	}

	wg.Wait()

	render(scrapers, os.Stdout)
}
