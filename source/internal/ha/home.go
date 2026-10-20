package ha

import (
	"cmp"
	"slices"
)

// Tile is one block of the panel's home page: an entity the panel controls,
// or what a device that only measures reads.
type Tile struct {
	// Entity is the entity the tile controls, or the sensor device's most
	// telling reading.
	Entity State
	// Readings are the sensor device's other readings, the most telling
	// first.
	Readings []State
	// Name is the sensor device's name, else the entity's.
	Name string
}

// Section holds the tiles of one area, named as Home Assistant names it,
// with the area's ID. The tiles of no area come last, in a section of no
// name and no ID.
type Section struct {
	ID    string
	Area  string
	Tiles []Tile
}

// tiers rank the domains the panel controls. A device shows its controls of
// the first tier it has: an air conditioner's climate entity, not the
// switch of its display; a power strip's switches.
var tiers = map[string]int{
	"climate": 1, "light": 1, "fan": 1, "humidifier": 1, "water_heater": 1,
	"cover": 1, "valve": 1, "lock": 1, "media_player": 1, "vacuum": 1,
	"lawn_mower": 1, "alarm_control_panel": 1, "camera": 1,
	"switch": 2, "siren": 2,
	"select": 3, "button": 3,
}

// order is the order of the tiles within an area, by domain; a sensor
// device's come last.
var order = []string{
	"climate", "light", "switch", "fan", "cover", "humidifier", "water_heater",
	"media_player", "vacuum", "lawn_mower", "lock", "valve", "alarm_control_panel",
	"siren", "camera", "select", "button",
}

// Readings of a sensor device, the most telling first. A device whose
// sensors read none of these, as a phone's, has no tile.
var (
	binaryOrder = []string{
		"door", "garage_door", "window", "opening", "lock", "motion", "occupancy",
		"presence", "moisture", "smoke", "gas", "carbon_monoxide", "safety",
		"vibration", "sound", "heat", "cold", "light", "plug", "power", "moving",
		"running", "",
	}
	sensorOrder = []string{
		"temperature", "humidity", "pm25", "pm10", "pm1", "carbon_dioxide", "aqi",
		"volatile_organic_compounds", "volatile_organic_compounds_parts",
		"carbon_monoxide", "nitrogen_dioxide", "ozone", "illuminance", "pressure",
		"atmospheric_pressure", "moisture", "power", "energy", "gas", "water",
	}
)

// readingRank orders a sensor device's readings, or is -1 for one the tile
// leaves out, as a battery level.
func readingRank(s State) int {
	class := s.Str("device_class")
	if s.Domain() == "binary_sensor" {
		return slices.Index(binaryOrder, class)
	}
	if class == "" && (s.Str("unit_of_measurement") == "°C" || s.Str("unit_of_measurement") == "°F") {
		class = "temperature"
	}
	if i := slices.Index(sensorOrder, class); i >= 0 {
		return len(binaryOrder) + i
	}
	return -1
}

// Arrange picks the tiles of the home page from every entity's state and
// sorts them into areas. Each device shows its controls, or else one tile of
// its readings; reg also leaves out what Home Assistant's dashboards hide.
// Without reg, which tells devices apart, each switch and main control
// shows alone and readings do not.
func Arrange(states []State, reg *Registry) []Section {
	type candidate struct {
		state  State
		index  int
		device string // ID, "" for no device
		owner  string // the device's name, else the entity's
		area   string // ID
	}
	var cands []candidate
	devices := map[string]Device{}
	for i, s := range states {
		d := s.Domain()
		if _, ok := tiers[d]; !ok && d != "sensor" && d != "binary_sensor" {
			continue
		}
		c := candidate{state: s, index: i, owner: s.Name()}
		if reg != nil {
			e := reg.Entities[s.EntityID]
			if e.Hidden {
				continue
			}
			if dev, ok := reg.Devices[e.Device]; ok {
				if dev.Service {
					continue
				}
				c.device, c.area, devices[e.Device] = e.Device, dev.Area, dev
				if dev.Name != "" {
					c.owner = dev.Name
				}
			}
			if e.Area != "" {
				c.area = e.Area
			}
		}
		cands = append(cands, c)
	}
	// The first tier of controls each device has.
	best := map[string]int{}
	for _, c := range cands {
		if t, ok := tiers[c.state.Domain()]; ok && c.device != "" {
			if b, seen := best[c.device]; !seen || t < b {
				best[c.device] = t
			}
		}
	}
	type placed struct {
		tile  Tile
		area  string
		rank  int
		owner string
		index int
	}
	var tiles []placed
	readings := map[string][]candidate{}
	var measured []string // devices with readings, in Home Assistant's order
	for _, c := range cands {
		d := c.state.Domain()
		if t, ok := tiers[d]; ok {
			if c.device != "" && t != best[c.device] || reg == nil && t > 2 {
				continue
			}
			tiles = append(tiles, placed{Tile{Entity: c.state, Name: c.state.Name()}, c.area, slices.Index(order, d), c.owner, c.index})
			continue
		}
		if _, controls := best[c.device]; c.device == "" || controls || readingRank(c.state) < 0 {
			continue
		}
		if readings[c.device] == nil {
			measured = append(measured, c.device)
		}
		readings[c.device] = append(readings[c.device], c)
	}
	for _, id := range measured {
		rs := readings[id]
		slices.SortStableFunc(rs, func(x, y candidate) int { return cmp.Compare(readingRank(x.state), readingRank(y.state)) })
		t := Tile{Entity: rs[0].state, Name: rs[0].owner}
		for _, r := range rs[1:] {
			t.Readings = append(t.Readings, r.state)
		}
		area := devices[id].Area
		if rs[0].area != "" {
			area = rs[0].area
		}
		tiles = append(tiles, placed{t, area, len(order), rs[0].owner, rs[0].index})
	}
	slices.SortStableFunc(tiles, func(x, y placed) int {
		return cmp.Or(cmp.Compare(x.rank, y.rank), cmp.Compare(x.owner, y.owner), cmp.Compare(x.index, y.index))
	})

	var sections []Section
	known := map[string]bool{}
	if reg != nil {
		for _, a := range reg.Areas {
			known[a.ID] = true
			var s Section
			for _, t := range tiles {
				if t.area == a.ID {
					s.Tiles = append(s.Tiles, t.tile)
				}
			}
			if len(s.Tiles) > 0 {
				s.ID, s.Area = a.ID, a.Name
				sections = append(sections, s)
			}
		}
	}
	var rest Section
	for _, t := range tiles {
		if !known[t.area] {
			rest.Tiles = append(rest.Tiles, t.tile)
		}
	}
	if len(rest.Tiles) > 0 {
		sections = append(sections, rest)
	}
	return sections
}
