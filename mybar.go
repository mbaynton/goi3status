// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// sample-bar demonstrates a sample i3bar built using barista.
package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/soumya92/barista/modules/funcs"

	"github.com/soumya92/barista/modules/diskspace"

	"github.com/soumya92/barista/modules/wlan"

	"github.com/soumya92/barista/modules/battery"

	"github.com/soumya92/barista/bar"
	"github.com/soumya92/barista/colors"
	"github.com/soumya92/barista/modules/clock"
	"github.com/soumya92/barista/modules/group"
	"github.com/soumya92/barista/modules/meminfo"
	"github.com/soumya92/barista/modules/netspeed"
	"github.com/soumya92/barista/modules/sysinfo"
	"github.com/soumya92/barista/modules/weather"
	"github.com/soumya92/barista/modules/weather/openweathermap"
	"github.com/soumya92/barista/outputs"
	"github.com/soumya92/barista/pango"
	"github.com/soumya92/barista/pango/icons/fontawesome"
	"github.com/soumya92/barista/pango/icons/material"
	"github.com/soumya92/barista/pango/icons/typicons"
)

var spacer = pango.Span(" ", pango.XXSmall)

func truncate(in string, l int) string {
	if len([]rune(in)) <= l {
		return in
	}
	return string([]rune(in)[:l-1]) + "⋯"
}

func hms(d time.Duration) (h int, m int, s int) {
	h = int(d.Hours())
	m = int(d.Minutes()) % 60
	s = int(d.Seconds()) % 60
	return
}

func formatMediaTime(d time.Duration) string {
	h, m, s := hms(d)
	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%d:%02d", m, s)
}

/*func mediaFormatFunc(m media.Info) bar.Output {
	if m.PlaybackStatus == media.Stopped || m.PlaybackStatus == media.Disconnected {
		return nil
	}
	artist := truncate(m.Artist, 20)
	title := truncate(m.Title, 40-len(artist))
	if len(title) < 20 {
		artist = truncate(m.Artist, 40-len(title))
	}
	var iconAndPosition pango.Node
	if m.PlaybackStatus == media.Playing {
		iconAndPosition = pango.Span(
			colors.Hex("#f70"),
			fontawesome.Icon("music"),
			spacer,
			formatMediaTime(m.Position()),
			"/",
			formatMediaTime(m.Length),
		)
	} else {
		iconAndPosition = fontawesome.Icon("music", colors.Hex("#f70"))
	}
	return outputs.Pango(iconAndPosition, spacer, title, " - ", artist)
}*/

func startTaskManager(e bar.Event) {
	if e.Button == bar.ButtonLeft {
		exec.Command("xfce4-taskmanager").Run()
	}
}

func home(path string) string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return filepath.Join(usr.HomeDir, path)
}

func main() {
	material.Load("/usr/share/baristafonts-full/material-design-icons")
	// materialCommunity.Load(home(".fonts/Github/MaterialDesign-Webfont"))
	typicons.Load("/usr/share/baristafonts-full/typicons.font")
	// ionicons.Load(home(".fonts/Github/ionicons"))
	fontawesome.Load("/usr/share/baristafonts-full/Font-Awesome")

	colors.LoadFromMap(map[string]string{
		"good":     "#6d6",
		"degraded": "#dd6",
		"bad":      "#d66",
		"dim-icon": "#777",
	})

	userhint := funcs.Once(func(m funcs.Module) {
		var icon string
		user, _ := user.Current()
		if user.Uid == "19884" {
			icon = "graduation-cap"
		} else {
			icon = "child"
		}
		m.Output(outputs.Pango(
			fontawesome.Icon(icon),
		))
	})

	localtime := clock.New().OutputFunc(func(now time.Time) bar.Output {
		return outputs.Pango(
			material.Icon("today", colors.Scheme("dim-icon")),
			now.Format("Mon Jan 2 "),
			now.Format("3:04:05 PM"),
		)
	}).OnClick(func(e bar.Event) {
		if e.Button == bar.ButtonLeft {
			exec.Command("gsimplecal").Run()
		}
	})

	// Weather information comes from OpenWeatherMap.
	// https://openweathermap.org/api.
	wthr := weather.New(
		openweathermap.Zipcode("55316", "US").Build(),
	).OutputFunc(func(w weather.Weather) bar.Output {
		iconName := ""
		switch w.Condition {
		case weather.Thunderstorm,
			weather.TropicalStorm,
			weather.Hurricane:
			iconName = "stormy"
		case weather.Drizzle,
			weather.Hail:
			iconName = "shower"
		case weather.Rain:
			iconName = "downpour"
		case weather.Snow,
			weather.Sleet:
			iconName = "snow"
		case weather.Mist,
			weather.Smoke,
			weather.Whirls,
			weather.Haze,
			weather.Fog:
			iconName = "windy-cloudy"
		case weather.Clear:
			if !w.Sunset.IsZero() && time.Now().After(w.Sunset) {
				iconName = "night"
			} else {
				iconName = "sunny"
			}
		case weather.PartlyCloudy:
			iconName = "partly-sunny"
		case weather.Cloudy, weather.Overcast:
			iconName = "cloudy"
		case weather.Tornado,
			weather.Windy:
			iconName = "windy"
		}
		if iconName == "" {
			iconName = "warning-outline"
		} else {
			iconName = "weather-" + iconName
		}
		return outputs.Pango(
			// typicons.Icon(iconName), spacer,
			pango.Textf("%d°F", w.Temperature.F()),
		)
	}).OnClick(func(_ weather.Weather, e bar.Event) {
		if e.Button == bar.ButtonLeft {
			exec.Command("firefox", "https://www.accuweather.com/en/us/champlin-mn/55316/hourly-weather-forecast/333873").Run()
		}
	})

	/*vol := volume.Mixer("usb-Lenovo_ThinkPad_Thunderbolt_3_Dock_USB_Audio_000000000000-00-USB", "Master").OutputFunc(func(v volume.Volume) bar.Output {
		if v.Mute {
			return outputs.
				Pango(ionicons.Icon("volume-mute"), "MUT").
				Color(colors.Scheme("degraded"))
		}
		iconName := "low"
		pct := v.Pct()
		if pct > 66 {
			iconName = "high"
		} else if pct > 33 {
			iconName = "medium"
		}
		return outputs.Pango(
			ionicons.Icon("volume-"+iconName),
			spacer,
			pango.Textf("%2d%%", pct),
		)
	})*/

	loadAvg := sysinfo.New().OutputFunc(func(s sysinfo.Info) bar.Output {
		out := outputs.Textf("%0.2f %0.2f", s.Loads[0], s.Loads[2])
		// Load averages are unusually high for a few minutes after boot.
		if s.Uptime < 10*time.Minute {
			// so don't add colours until 10 minutes after system start.
			return out
		}
		switch {
		case s.Loads[0] > 7, s.Loads[2] > 64:
			out.Urgent(true)
		case s.Loads[0] > 3, s.Loads[2] > 2.25:
			out.Color(colors.Scheme("bad"))
		case s.Loads[0] > 2, s.Loads[2] > 2:
			out.Color(colors.Scheme("degraded"))
		}
		return out
	}).OnClick(startTaskManager)

	freeMem := meminfo.New().OutputFunc(func(m meminfo.Info) bar.Output {
		out := outputs.Pango(material.Icon("memory"), m.Available().IEC())
		freeGigs := m.Available().In("GiB")
		switch {
		case freeGigs < 0.5:
			out.Urgent(true)
		case freeGigs < 1:
			out.Color(colors.Scheme("bad"))
		case freeGigs < 2:
			out.Color(colors.Scheme("degraded"))
		case freeGigs > 6:
			out.Color(colors.Scheme("good"))
		}
		return out
	}).OnClick(startTaskManager)

	wifi_iface := "wlp4s0"
	net := netspeed.New(wifi_iface).
		RefreshInterval(2 * time.Second).
		OutputFunc(func(s netspeed.Speeds) bar.Output {
			return outputs.Pango(
				fontawesome.Icon("upload"), spacer, pango.Textf("%5s", s.Tx.SI()),
				pango.Span(" ", pango.Small),
				fontawesome.Icon("download"), spacer, pango.Textf("%5s", s.Rx.SI()),
			)
		})

	// bat := battery.Default().OutputTemplate(outputs.TextTemplate(`BATT {{.RemainingPct}}% {{.RemainingTime}}`))
	batBasics := battery.Default().OutputFunc(func(b battery.Info) bar.Output {
		var icon string
		var color bar.Color
		remaining := b.Remaining()
		if remaining >= .683333 {
			icon = "battery-full"
		} else if remaining >= .366666 {
			icon = "battery-half"
			color = colors.Scheme("degraded")
		} else if remaining >= .5 {
			icon = "battery-quarter"
			color = colors.Scheme("bad")
		} else {
			icon = "battery-empty"
		}

		if b.PluggedIn() {
			color = colors.Scheme("good")
		}

		out := outputs.Pango(
			fontawesome.Icon(icon),
			spacer,
			pango.Textf("%.1f%%/%.1f W", b.Remaining()*100, b.Power),
		)

		out.Color(color)
		if remaining < .05 {
			out.Urgent(true)
		}
		return out
	})

	batTime := battery.Default().OutputTemplate(outputs.TextTemplate(`{{.RemainingTime}}`))

	// wifi := wlan.New(wifi_iface).OutputTemplate(outputs.TextTemplate("W: {{if .Connected}}{{.SSID}}{{else}}None{{end}}"))
	wifi := wlan.New(wifi_iface).OutputFunc(func(w wlan.Info) bar.Output {
		var items []interface{}
		var icon string
		if w.Connected() {
			icon = "network-wifi"
		} else {
			icon = "signal-wifi-off"
		}

		items = append(items, material.Icon(icon))

		if w.Connected() {
			items = append(items, spacer, pango.Textf("%s", w.SSID))
		}
		return outputs.Pango(items...)
	})

	diskOs := diskspace.New("/").OutputTemplate(outputs.TextTemplate(`/: {{.Available.In "GiB" | printf "%.1f"}} GB`))
	diskHome := diskspace.New("/home").OutputTemplate(outputs.TextTemplate(`/home: {{.Available.In "GiB" | printf "%.1f"}} GB`))

	thirtySeconds, _ := time.ParseDuration("30s")
	backlight := funcs.Every(thirtySeconds, func(m funcs.Module) {
		brightnessBytes, err := ioutil.ReadFile("/sys/class/backlight/intel_backlight/brightness")
		if err != nil {
			m.Error(err)
			return
		}
		level, err := strconv.ParseFloat(
			strings.Trim(string(brightnessBytes), "\n "),
			32,
		)
		if err != nil {
			m.Error(err)
			return
		}
		backlightPct := (level / 1060.0) * 100

		m.Output(outputs.Pango(
			material.Icon("wb-incandescent"),
			pango.Textf("%.0f%%", backlightPct),
		))
	})

	g := group.Collapsing()

	batTimeGroup := group.Collapsing()
	batTimeGroup.Collapse()

	panic(bar.Run(
		backlight,
		wifi,
		g.Add(net),
		g.Add(diskOs),
		g.Add(diskHome),
		g.Add(freeMem),
		g.Button(outputs.Text("+"), outputs.Text("-")),
		wthr,
		batBasics,
		batTimeGroup.Add(batTime),
		batTimeGroup.Button(outputs.Text("*"), outputs.Text("-")),
		loadAvg,
		localtime,
		userhint,
	))
}
