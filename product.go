package coffeemenu

import (
	"strings"
)

type Product struct {
	Name string
	Url  string
}

type Products []Product

func (p Products) Len() int           { return len(p) }
func (p Products) Less(i, j int) bool { return strings.Compare(p[i].Name, p[j].Name) < 0 }
func (p Products) Swap(i, j int)      { (p)[i], (p)[j] = (p)[j], (p)[i] }
