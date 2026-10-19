package ha

import (
	"math"
	"slices"
)

// Color is a color a light shines in: a white of a color temperature, or a
// hue and saturation.
type Color struct {
	// Kelvin is a white's color temperature, else 0.
	Kelvin float64
	// Hue, 0 to 360, and Saturation, 0 to 100, are a color's.
	Hue, Saturation float64
	// Name says which color it is, as "red", for the panel to name; whites
	// and colors chosen freely have none.
	Name string
}

// White is the white of color temperature k.
func White(k float64) Color { return Color{Kelvin: math.Round(k)} }

// Hues are the colors a light that takes colors offers one click away.
var Hues = []Color{
	{Hue: 0, Saturation: 80, Name: "red"},
	{Hue: 28, Saturation: 85, Name: "orange"},
	{Hue: 50, Saturation: 80, Name: "yellow"},
	{Hue: 125, Saturation: 70, Name: "green"},
	{Hue: 215, Saturation: 80, Name: "blue"},
	{Hue: 280, Saturation: 65, Name: "purple"},
}

// The bounds of a light's color temperatures when it does not say them.
const warmest, coolest = 2000, 6500

// Whites reports the color temperatures, in kelvin, the light shines in,
// when it shines in whites of more than one.
func (s State) Whites() (warm, cool float64, ok bool) {
	if !slices.Contains(s.Strs("supported_color_modes"), "color_temp") {
		return 0, 0, false
	}
	warm, cool = warmest, coolest
	if v, has := s.Num("min_color_temp_kelvin"); has && v > 0 {
		warm = v
	}
	if v, has := s.Num("max_color_temp_kelvin"); has && v > warm {
		cool = v
	}
	return warm, cool, cool > warm
}

// Colorful reports whether the light shines in colors other than whites.
func (s State) Colorful() bool {
	return slices.ContainsFunc(s.Strs("supported_color_modes"), func(m string) bool {
		return m == "hs" || m == "xy" || m == "rgb" || m == "rgbw" || m == "rgbww"
	})
}

// Swatches are the whites a light offers one click away, up to four from
// its warmest to its coolest, evenly apart as the eye sees them.
func (s State) Swatches() []Color {
	warm, cool, ok := s.Whites()
	if !ok {
		return nil
	}
	const n = 4
	var out []Color
	for i := range n {
		// Mireds, a million over kelvin, space whites as the eye does.
		m := 1e6/warm + (1e6/cool-1e6/warm)*float64(i)/(n-1)
		w := White(max(warm, min(cool, math.Round(1e6/m/100)*100)))
		if len(out) == 0 || out[len(out)-1] != w {
			out = append(out, w)
		}
	}
	return out
}

// Shining is the color the light shines in now, if it is on and says.
func (s State) Shining() (Color, bool) {
	if !s.Active() {
		return Color{}, false
	}
	if s.Str("color_mode") == "color_temp" {
		k, ok := s.Num("color_temp_kelvin")
		return White(k), ok && k > 0
	}
	hs, _ := s.Attributes["hs_color"].([]any)
	if len(hs) != 2 {
		return Color{}, false
	}
	h, ok1 := hs[0].(float64)
	sat, ok2 := hs[1].(float64)
	return Color{Hue: h, Saturation: sat}, ok1 && ok2
}

// Near reports whether c and d are the same color to the eye, as a color
// Home Assistant reports back after converting the one sent.
func (c Color) Near(d Color) bool {
	if c.Kelvin > 0 || d.Kelvin > 0 {
		return c.Kelvin > 0 && d.Kelvin > 0 && math.Abs(1e6/c.Kelvin-1e6/d.Kelvin) < 6
	}
	dh := math.Abs(c.Hue - d.Hue)
	dh = min(dh, 360-dh)
	return dh < 6 && math.Abs(c.Saturation-d.Saturation) < 6
}

// Set is the action that turns the light on in c, showing c until Home
// Assistant reports it.
func (c Color) Set() Action {
	r, g, b := c.RGB()
	rgb := []any{float64(r), float64(g), float64(b)}
	if c.Kelvin > 0 {
		return Action{Name: "turn_on", Service: "turn_on", Data: map[string]any{"color_temp_kelvin": c.Kelvin}, Next: "on",
			Shows: map[string]any{"color_mode": "color_temp", "color_temp_kelvin": c.Kelvin, "rgb_color": rgb}}
	}
	hs := []float64{math.Round(c.Hue*10) / 10, math.Round(c.Saturation*10) / 10}
	return Action{Name: "turn_on", Service: "turn_on", Data: map[string]any{"hs_color": hs}, Next: "on",
		Shows: map[string]any{"color_mode": "hs", "hs_color": []any{hs[0], hs[1]}, "rgb_color": rgb}}
}

// RGB is how c looks on a screen, at full brightness.
func (c Color) RGB() (r, g, b uint8) {
	if c.Kelvin > 0 {
		return kelvinRGB(c.Kelvin)
	}
	h := math.Mod(c.Hue, 360) / 60
	sat := max(0, min(1, c.Saturation/100))
	f := h - math.Floor(h)
	p, q, t := 1-sat, 1-sat*f, 1-sat*(1-f)
	var rf, gf, bf float64
	switch int(h) {
	case 0:
		rf, gf, bf = 1, t, p
	case 1:
		rf, gf, bf = q, 1, p
	case 2:
		rf, gf, bf = p, 1, t
	case 3:
		rf, gf, bf = p, q, 1
	case 4:
		rf, gf, bf = t, p, 1
	default:
		rf, gf, bf = 1, p, q
	}
	return byte255(rf * 255), byte255(gf * 255), byte255(bf * 255)
}

// kelvinRGB approximates the color of a black body at k kelvin, as Tanner
// Helland fit it.
func kelvinRGB(k float64) (r, g, b uint8) {
	t := max(1000, min(40000, k)) / 100
	rf, gf, bf := 255.0, 255.0, 255.0
	if t <= 66 {
		gf = 99.4708025861*math.Log(t) - 161.1195681661
		if t <= 19 {
			bf = 0
		} else {
			bf = 138.5177312231*math.Log(t-10) - 305.0447927307
		}
	} else {
		rf = 329.698727446 * math.Pow(t-60, -0.1332047592)
		gf = 288.1221695283 * math.Pow(t-60, -0.0755148492)
	}
	return byte255(rf), byte255(gf), byte255(bf)
}

func byte255(v float64) uint8 { return uint8(max(0, min(255, math.Round(v)))) }

// Stops are the levels the panel offers one click away below the slider:
// a cover's or a valve's quarters.
func (l Level) Stops() []float64 {
	if l.Name != "position" {
		return nil
	}
	return []float64{0, 25, 50, 75, 100}
}
