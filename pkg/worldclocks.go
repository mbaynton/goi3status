package pkg

import (
	"barista.run/bar"
	"barista.run/modules/clock"
	"barista.run/outputs"
	"barista.run/pango"
	"fmt"
	"time"
)

var spacer = pango.Text(" ").XXSmall()

func makeTzClock(lbl, tzName string, useMilFormat bool) bar.Module {
	homeTz, err := time.LoadLocation("America/Chicago")
	if err != nil {
		homeTz = time.UTC
	}

	clockTz, err := time.LoadLocation(tzName)
	if err != nil {
		panic(err)
	}
	c := clock.Zone(clockTz)
	return c.Output(time.Minute, func(now time.Time) bar.Output {
		homeTime := now.In(homeTz)
		clockTime := now.In(clockTz)
		homeHours := time.Date(homeTime.Year(), homeTime.Month(), homeTime.Day(), homeTime.Hour(), homeTime.Minute(), 0, 0, time.UTC)
		clockHours := time.Date(clockTime.Year(), clockTime.Month(), clockTime.Day(), clockTime.Hour(), clockTime.Minute(), 0, 0, time.UTC)
		homeDiff := clockHours.Sub(homeHours).Hours()
		format := "03:04 PM"
		if useMilFormat {
			format = "15:04"
		}
		return outputs.Pango(pango.Text(lbl).Smaller(), spacer, now.Format(format), spacer, pango.Text(fmt.Sprintf("(%+.0fh)", homeDiff)))
	})
}

func GetWorldClocks() []bar.Module {
	clocks := []bar.Module{
		makeTzClock("West Coast", "America/Los_Angeles", false),
		makeTzClock("Central", "America/Chicago", false),
		makeTzClock("East Coast", "America/New_York", false),
		makeTzClock("UTC", "Etc/UTC", true),
	}

	return clocks
}
