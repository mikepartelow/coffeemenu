package main

import (
	"fmt"
	"sort"
	"sync"

	"github.com/mikepartelow/coffeemenu"
	"github.com/mikepartelow/coffeemenu/internal/lineacaffe"
	"github.com/mikepartelow/coffeemenu/internal/mrespresso"
	"github.com/mikepartelow/coffeemenu/internal/rebootroasting"
	"github.com/rs/zerolog/log"
)

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
				log.Error().Msgf("Error while scraping %q: %v", s.GetName(), err)
			}
			wg.Done()
		}()
	}

	wg.Wait()

	for _, s := range scrapers {
		products := s.GetProducts()
		sort.Sort(products)

		fmt.Println("# ", s.GetName())
		for _, p := range products {
			fmt.Println("- ", p.Name)
		}
		fmt.Println("> ", s.GetStats())
		fmt.Println()
	}
}
