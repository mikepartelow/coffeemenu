package main

import (
	"os"
	"sync"

	_ "embed"

	"github.com/mikepartelow/coffeemenu"
	"github.com/rs/zerolog/log"
)

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

	coffeemenu.Render(scrapers, os.Stdout)
}
