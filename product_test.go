package coffeemenu_test

import (
	"sort"
	"testing"

	"github.com/mikepartelow/coffeemenu"
	"github.com/stretchr/testify/assert"
)

func TestSort(t *testing.T) {
	got := coffeemenu.Products{
		{
			Name: "Zork Beanery",
			Url:  "http://zork.com",
		},
		{
			Name: "Alpha Beta",
			Url:  "http://alphabeta.com",
		},
		{
			Name: "Morpheus Coffee",
			Url:  "http://morpheus-coffee.com",
		},
	}

	want := []string{"Alpha Beta", "Morpheus Coffee", "Zork Beanery"}

	sort.Sort(got)

	for i := range got {
		assert.Equal(t, got[i].Name, want[i])
	}
}
