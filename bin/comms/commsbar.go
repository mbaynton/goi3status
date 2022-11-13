package main

import (
	barista "barista.run"
	"github.com/mbaynton/goi3status/pkg"
)

func main() {
	worldClocks := pkg.GetWorldClocks()
	for _, clock := range worldClocks {
		barista.Add(clock)
	}

	panic(barista.Run())
}
