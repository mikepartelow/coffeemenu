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

func main() {
	boring := flag.Bool("boring", false, "render boring MarkDown")
	csv := flag.Bool("csv", false, "render CSV")
	flag.Parse()

	var scrapers []*coffeemenu.Scraper

	if args := flag.Args(); len(args) == 1 {
		scrapers = append(scrapers, readScraper(args[0]))
	} else {
		scrapers = readScrapers()
	}

	out := render(scrape(scrapers), *csv, *boring)

	fmt.Println(out)
}

//go:embed sites/*.json
var sitesFS embed.FS

func readScraper(name string) *coffeemenu.Scraper {
	site, err := coffeemenu.ReadSite(sitesFS, name)
	if err != nil {
		log.Fatal().Err(err).Msgf("Couldn't read site %q.", name)
	}

	return coffeemenu.NewScraper(*site, nil)
}

func readScrapers() []*coffeemenu.Scraper {
	sites, err := coffeemenu.ReadSites(sitesFS)
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't read sites.")
	}

	var scrapers []*coffeemenu.Scraper
	for _, site := range sites {
		scrapers = append(scrapers, coffeemenu.NewScraper(site, nil))
	}

	return scrapers
}

func scrape(scrapers []*coffeemenu.Scraper) []*coffeemenu.Scraper {
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

	return scrapers
}

func render(scrapers []*coffeemenu.Scraper, csv, boring bool) string {
	if csv {
		return coffeemenu.CSV(scrapers)
	}

	var buf bytes.Buffer
	coffeemenu.Render(scrapers, &buf)

	if boring {
		return buf.String()
	}

	return coffeemenu.Glamourize(buf.String())
}
