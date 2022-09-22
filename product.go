package coffeemenu

import (
	"strings"

	"github.com/savioxavier/termlink"
)

type Product struct {
	Name string
	Url  string
}

type Products []Product

func (p Products) Len() int           { return len(p) }
func (p Products) Less(i, j int) bool { return strings.Compare(p[i].Name, p[j].Name) < 0 }
func (p Products) Swap(i, j int)      { (p)[i], (p)[j] = (p)[j], (p)[i] }

func (p Product) TermLink() string {
	return termlink.Link(p.Name, p.Url)
}
