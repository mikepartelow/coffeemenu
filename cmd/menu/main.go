package main

import (
	"bytes"
	"flag"
	"fmt"
	"sync"

	_ "embed"

	"github.com/mikepartelow/coffeemenu"
	"github.com/rs/zerolog/log"
)

func main() {
	boring := flag.Bool("boring", false, "render boring MarkDown")
	flag.Parse()

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

	var out string
	var buf bytes.Buffer
	coffeemenu.Render(scrapers, &buf)

	if *boring {
		out = buf.String()
	} else {
		out = coffeemenu.Glamourize(buf.String())
	}
	fmt.Println(out)
}
