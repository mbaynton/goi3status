package main

import (
	barista "barista.run"
	"barista.run/modules/static"
	"barista.run/outputs"
	"fmt"
	"github.com/mbaynton/goi3status/pkg"
	"os"
	"time"
)

func main() {
	barista.Add(static.New(outputs.Text("World Clocks: ")))
	
	worldClocks := pkg.GetWorldClocks()
	for _, clock := range worldClocks {
		barista.Add(clock)
	}

	err := barista.Run()
	if err != nil {
		os.WriteFile("/tmp/commsbar-exit-log", []byte(fmt.Sprintf("%v: %+v\n", time.Now(), err)), os.ModeAppend)
		os.Exit(1)
	}
	os.Exit(0)
}
