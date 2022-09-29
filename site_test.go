package coffeemenu_test

import (
	"testing"
	"testing/fstest"

	"github.com/mikepartelow/coffeemenu"
	"github.com/stretchr/testify/assert"
)

var testfs fstest.MapFS = fstest.MapFS{
	"sites/bananapantsroasting.json": {
		Data: []byte(`{ "name": "Banana Pants Roasting",
    					"urls": [ "https://www.bananapantsroasting.com/buystuff" ],
    					"scrapespec": {
        					"container": "div.product",
        					"name": [ "What", "Is", "This" ],
        					"url": [ "Whats", "ThisNow" ] }
					}`),
	},
}

func TestReadSites(t *testing.T) {
	sites, err := coffeemenu.ReadSites(testfs)
	assert.NoError(t, err)
	assert.Len(t, sites, 1)
	assert.Equal(t, sites[0].Name, "Banana Pants Roasting")
	assert.Equal(t, sites[0].Urls, []string{"https://www.bananapantsroasting.com/buystuff"})
	assert.Equal(t, sites[0].ScrapeSpec.Name, []string{"What", "Is", "This"})
}

func TestReadSite(t *testing.T) {
	site, err := coffeemenu.ReadSite(testfs, "bananapantsroasting")
	assert.NoError(t, err)
	assert.Equal(t, site.Name, "Banana Pants Roasting")
	assert.Equal(t, site.Urls, []string{"https://www.bananapantsroasting.com/buystuff"})
	assert.Equal(t, site.ScrapeSpec.Name, []string{"What", "Is", "This"})

}
