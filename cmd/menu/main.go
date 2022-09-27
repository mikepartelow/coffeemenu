package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"sync"

	_ "embed"

	"github.com/mikepartelow/coffeemenu"
	"github.com/rs/zerolog/log"
)

//go:embed sites/*.json
var siteFS embed.FS

func main() {
	boring := flag.Bool("boring", false, "render boring MarkDown")
	csv := flag.Bool("csv", false, "render CSV")
	flag.Parse()

	sites, err := coffeemenu.ReadSites(siteFS)
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't read sites.")
	}

	var scrapers []*coffeemenu.Scraper
	for _, site := range sites {
		scrapers = append(scrapers, coffeemenu.NewScraper(site, nil))
	}

	var wg sync.WaitGroup

	for _, s := range scrapers {
		wg.Add(1)
		go func(s *coffeemenu.Scraper) {
			defer wg.Done()
			if err := s.Scrape(); err != nil {
				log.Error().Err(err).Msgf("Error while scraping %q", s.Name())
			}
		}(s)
	}

	wg.Wait()

	var out string

	if *csv {
		out = coffeemenu.CSV(scrapers)
	} else {
		var buf bytes.Buffer
		coffeemenu.Render(scrapers, &buf)

		if *boring {
			out = buf.String()
		} else {
			out = coffeemenu.Glamourize(buf.String())
		}
	}
	fmt.Println(out)
}
