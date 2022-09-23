package main

import (
	"bytes"
	"fmt"
	"sync"

	_ "embed"

	"github.com/charmbracelet/glamour"
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

	r, err := glamour.NewTermRenderer(
		glamour.WithStylePath("dracula"),
		glamour.WithWordWrap(0),
	)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	var buf bytes.Buffer
	coffeemenu.Render(scrapers, &buf)

	out, err := r.Render(buf.String())
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	fmt.Println(out)
}
