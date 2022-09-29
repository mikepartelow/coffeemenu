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
	assert.Equal(t, "bananapantsroasting", sites[0].ID)
}

func TestReadSite(t *testing.T) {
	site, err := coffeemenu.ReadSite(testfs, "bananapantsroasting")
	assert.NoError(t, err)
	assert.Equal(t, "Banana Pants Roasting", site.Name)
	assert.Equal(t, []string{"https://www.bananapantsroasting.com/buystuff"}, site.Urls)
	assert.Equal(t, []string{"What", "Is", "This"}, site.ScrapeSpec.Name)
	assert.Equal(t, "bananapantsroasting", site.ID)
}
