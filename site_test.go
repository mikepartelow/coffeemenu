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

func TestReadScrapers(t *testing.T) {
	scrapers, err := coffeemenu.ReadScrapers(testfs)
	assert.NoError(t, err)

	assert.Len(t, scrapers, 1)
	assert.Equal(t, scrapers[0].Name(), "Banana Pants Roasting")
}
