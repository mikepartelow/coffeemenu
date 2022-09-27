package coffeemenu

import (
	"encoding/json"
	"io/fs"
	"path"

	"github.com/rs/zerolog/log"
)

type ScrapeSpec struct {
	Container string   `json:"container"`
	Name      []string `json:"name"`
	Url       []string `json:"url"`
}

type Site struct {
	Name       string     `json:"name"`
	Urls       []string   `json:"urls"`
	ScrapeSpec ScrapeSpec `json:"scrapespec"`
}

type ScrapeSpecFs interface {
	fs.ReadDirFS
	fs.ReadFileFS
}

func ReadSites(sitesFS ScrapeSpecFs) ([]Site, error) {
	var sites []Site

	siteFiles, err := sitesFS.ReadDir("sites")
	if err != nil {
		return nil, err
	}

	for _, file := range siteFiles {
		bytes, err := sitesFS.ReadFile(path.Join("sites", file.Name()))
		if err != nil {
			return nil, err
		}
		var site Site
		if err := json.Unmarshal(bytes, &site); err != nil {
			log.Error().Err(err).Send()
			return nil, err
		}
		sites = append(sites, site)
	}

	return sites, nil
}
