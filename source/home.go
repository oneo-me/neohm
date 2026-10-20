package main

import (
	"cmp"
	"maps"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/egoist/mygo/plugins/glass"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/neohm/source/internal/ha"
	"github.com/oneo-me/neohm/source/locales"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

// tileHeight is the height of a tile, about the width of a third of the
// panel.
const tileHeight = 112

// pinnedID is the ID of the group of pinned tiles above Home's areas, which
// no area of Home Assistant's slugs takes.
const pinnedID = "@pinned"

// homePage is where the panel opens: Home Assistant's devices as tiles, by
// area. radius rounds the corners the tiles scroll within.
func (a *application) homePage(c *ui.Context, radius float32) {
	pageFrame(c, radius, locales.Title(), func() {
		pageHeader(c, locales.Title(), func() {
			a.settingsButton(c)
		})
	}, func(top float32) {
		switch {
		case !a.configured() || !a.listed:
			a.setupView(c, top)
		case len(a.home) == 0:
			ui.Column(c).Fill().Center().Gap(10).Padding(top, 24, 48, 24).Children(func() {
				components.Icon(c, lucide.House).Size(44, 44).TextColor(inkMuted)
				components.Text(c, locales.NoDevices()).FontSize(17).Bold().TextColor(ink)
				components.Text(c, locales.NoDevicesBody()).Wrap().TextAlign(ui.Center).FontSize(13).TextColor(inkMuted).MaxWidth(260)
			})
		default:
			scrollContent(c, top, radius, 18, func() {
				// A failure after the devices were read leaves them as last
				// read.
				if e := a.errorDisplay(); e != "" {
					ui.Row(c).FillWidth().Gap(8).Padding(0, 4).AlignItems(ui.Center).Children(func() {
						components.Icon(c, lucide.TriangleAlert).Size(16, 16).TextColor(danger).Shrink(0)
						components.Text(c, e).FontSize(12).Wrap().TextColor(inkMuted).Grow(1).MinWidth(0)
					})
				}
				// Releasing a drag anywhere among the areas keeps the order it
				// shows; Escape, or releasing it elsewhere, restores the one
				// before.
				list := ui.Column(c).FillWidth().Gap(18)
				if _, ok := ui.Drop[homeDrag](list); ok && a.reorder != nil {
					a.reorder.dropped = true
				}
				list.Children(func() {
					for _, s := range a.home {
						a.area(c, s)
					}
				})
			})
			a.settleReorder(c)
		}
	})
}

// area shows the tiles of area s under its heading, which folds them away.
// While another area is dragged over it, that area moves above or below
// it, by the half of it under the pointer. The pinned tiles' group neither
// folds nor moves, and stays above the areas.
func (a *application) area(c *ui.Context, s ha.Section) {
	col := ui.Column(c.Key("area:" + s.ID)).FillWidth().Gap(10)
	if r := a.reorder; r != nil && !r.stale && r.drag.entity == "" && r.drag.area != s.ID && s.ID != pinnedID {
		if _, y, over := col.PointerPosition(); over {
			after := y >= col.LayoutBounds().H/2
			r.move = func() { a.moveArea(r.drag.area, s.ID, after) }
		}
	}
	collapsed := slices.Contains(a.config.Collapsed, s.ID)
	col.Children(func() {
		switch {
		case s.ID == pinnedID:
			collapsed = false
			ui.Row(c).FillWidth().Gap(6).Padding(2, 4).AlignItems(ui.Center).Children(func() {
				areaHeading(c, locales.Pinned()).Grow(1).MinWidth(0).Margin(0)
				components.Icon(c, lucide.Pin).Size(14, 14).TextColor(inkMuted)
			})
		case s.Area != "":
			a.areaHeader(c, s, s.Area, collapsed)
		case len(a.home) > 1:
			a.areaHeader(c, s, locales.OtherDevices(), collapsed)
		default:
			collapsed = false
		}
		if collapsed {
			return
		}
		ui.Grid(c).Columns(3).Gap(8).Children(func() {
			for _, t := range s.Tiles {
				a.tile(c, t, s)
			}
		})
	})
}

// areaHeader names area s, and folds its tiles away or shows them again as
// it is clicked. Dragging it moves the area among the others.
func (a *application) areaHeader(c *ui.Context, s ha.Section, name string, collapsed bool) {
	h := ui.ButtonBase(c).Row().FillWidth().Gap(6).Padding(2, 4).AlignItems(ui.Center).Label(name).Expanded(!collapsed)
	if a.areaCount() > 1 {
		d := homeDrag{area: s.ID}
		h.Drag(d)
		if h.Dragging() {
			a.dragging(d)
			h.Opacity(0.4)
		}
	}
	turn := float32(0)
	if collapsed {
		turn = 1
	}
	turn = h.Animate("turn", turn, 160*time.Millisecond)
	h.Children(func() {
		areaHeading(c, name).Grow(1).MinWidth(0).Margin(0)
		if collapsed {
			components.Text(c, strconv.Itoa(len(s.Tiles))).FontSize(12).TextColor(inkMuted)
		}
		components.Icon(c, lucide.ChevronDown).Size(14, 14).TextColor(inkMuted).Rotate(-90 * turn)
	})
	if h.Clicked() {
		a.toggleArea(s.ID)
	}
}

// areaCount is how many areas Home shows, without the pinned tiles' group.
func (a *application) areaCount() int {
	if len(a.home) > 0 && a.home[0].ID == pinnedID {
		return len(a.home) - 1
	}
	return len(a.home)
}

// toggleArea folds away the tiles of area id, or shows them again.
func (a *application) toggleArea(id string) {
	if i := slices.Index(a.config.Collapsed, id); i >= 0 {
		a.config.Collapsed = slices.Delete(slices.Clone(a.config.Collapsed), i, i+1)
	} else {
		a.config.Collapsed = append(slices.Clone(a.config.Collapsed), id)
	}
	a.save()
}

// togglePin pins the tile of entity id above Home's areas, after the tiles
// pinned before, or returns it to its area.
func (a *application) togglePin(id string) {
	if i := slices.Index(a.config.Pinned, id); i >= 0 {
		a.config.Pinned = slices.Delete(slices.Clone(a.config.Pinned), i, i+1)
	} else {
		a.config.Pinned = append(slices.Clone(a.config.Pinned), id)
	}
	a.arrange()
	a.save()
}

// homeDrag is an area of Home being dragged to another place among the
// others, or, when entity is set, the tile of entity within its area.
type homeDrag struct{ area, entity string }

// reorder is a drag moving an area or a tile of Home, as it goes: the
// order saved as it started, which it returns to unless it drops among the
// areas, and the move the frame found the pointer asks for.
type reorder struct {
	drag          homeDrag
	areas, pinned []string
	tiles         map[string][]string
	move          func()
	// active reports that the source dragged this frame, moved that the
	// order changed and dropped that the drag dropped among the areas.
	// stale reports that the frame's pointer is placed on the layout from
	// before the last move.
	active, moved, dropped, stale bool
}

// dragging reports that the source of d drags it this frame, which starts
// a reorder.
func (a *application) dragging(d homeDrag) {
	if a.reorder == nil {
		a.reorder = &reorder{drag: d, areas: a.config.AreaOrder, pinned: a.config.Pinned, tiles: maps.Clone(a.config.TileOrder)}
	}
	a.reorder.active = true
}

// settleReorder moves what the frame's drag asks to, once Home is built,
// and ends a drag no longer dragging: it keeps the order shown if the drag
// dropped among the areas, and restores the one before otherwise.
func (a *application) settleReorder(c *ui.Context) {
	r := a.reorder
	if r == nil {
		return
	}
	r.stale = r.move != nil
	if r.move != nil {
		r.move()
		r.move = nil
		c.Invalidate()
	}
	if r.active {
		r.active = false
		return
	}
	a.reorder = nil
	if !r.moved {
		return
	}
	if r.dropped {
		a.save()
		return
	}
	a.config.AreaOrder, a.config.TileOrder, a.config.Pinned = r.areas, r.tiles, r.pinned
	a.arrange()
	c.Invalidate()
}

// moveArea moves area id above area target, or below it after.
func (a *application) moveArea(id, target string, after bool) {
	i := slices.IndexFunc(a.home, func(s ha.Section) bool { return s.ID == id })
	to := slices.IndexFunc(a.home, func(s ha.Section) bool { return s.ID == target })
	if i < 0 || to < 0 {
		return
	}
	if after {
		to++
	}
	if i < to {
		to--
	}
	if i == to {
		return
	}
	s := a.home[i]
	a.home = slices.Insert(slices.Delete(a.home, i, i+1), to, s)
	a.config.AreaOrder = nil
	for _, s := range a.home {
		if s.ID != pinnedID {
			a.config.AreaOrder = append(a.config.AreaOrder, s.ID)
		}
	}
	a.reorder.moved = true
}

// moveTile moves the tile of entity id into the place of target's in area,
// which moves the tiles between them by one.
func (a *application) moveTile(area, id, target string) {
	k := slices.IndexFunc(a.home, func(s ha.Section) bool { return s.ID == area })
	if k < 0 {
		return
	}
	tiles := a.home[k].Tiles
	i := slices.IndexFunc(tiles, func(t ha.Tile) bool { return t.Entity.EntityID == id })
	to := slices.IndexFunc(tiles, func(t ha.Tile) bool { return t.Entity.EntityID == target })
	if i < 0 || to < 0 || i == to {
		return
	}
	t := tiles[i]
	tiles = slices.Insert(slices.Delete(tiles, i, i+1), to, t)
	a.home[k].Tiles = tiles
	a.reorder.moved = true
	if area == pinnedID {
		a.config.Pinned = tileIDs(tiles)
		return
	}
	a.config.TileOrder = maps.Clone(a.config.TileOrder)
	if a.config.TileOrder == nil {
		a.config.TileOrder = map[string][]string{}
	}
	a.config.TileOrder[area] = tileIDs(tiles)
}

func tileIDs(tiles []ha.Tile) []string {
	ids := make([]string, len(tiles))
	for i, t := range tiles {
		ids[i] = t.Entity.EntityID
	}
	return ids
}

// arrange shows the areas Home Assistant arranges on Home as the user
// arranged them: the tiles pinned in a group above the areas, then the
// areas, and the tiles of each, in the order they were dragged into. An
// area left without tiles drops out.
func (a *application) arrange() {
	var pinned []ha.Tile
	var home []ha.Section
	for _, s := range a.sections {
		s.Tiles = slices.DeleteFunc(slices.Clone(s.Tiles), func(t ha.Tile) bool {
			if !slices.Contains(a.config.Pinned, t.Entity.EntityID) {
				return false
			}
			pinned = append(pinned, t)
			return true
		})
		if len(s.Tiles) > 0 {
			home = append(home, s)
		}
	}
	home = arrange(home, a.config.AreaOrder, a.config.TileOrder)
	if len(pinned) > 0 {
		home = slices.Insert(home, 0, ha.Section{ID: pinnedID, Tiles: pinned})
		arrange(home[:1], nil, map[string][]string{pinnedID: a.config.Pinned})
	}
	a.home = home
}

// arrange orders the areas of home by their IDs in areas, and the tiles of
// each by their entities in tiles. What these miss follows in the order it
// has.
func arrange(home []ha.Section, areas []string, tiles map[string][]string) []ha.Section {
	slices.SortStableFunc(home, func(x, y ha.Section) int {
		return cmp.Compare(rank(areas, x.ID), rank(areas, y.ID))
	})
	for _, s := range home {
		if order, ok := tiles[s.ID]; ok {
			slices.SortStableFunc(s.Tiles, func(x, y ha.Tile) int {
				return cmp.Compare(rank(order, x.Entity.EntityID), rank(order, y.Entity.EntityID))
			})
		}
	}
	return home
}

// rank is the place of id in order, past its end when order misses it.
func rank(order []string, id string) int {
	if i := slices.Index(order, id); i >= 0 {
		return i
	}
	return len(order)
}

// settingsButton opens Settings, marked with a dot while a newer version
// waits there to be installed.
func (a *application) settingsButton(c *ui.Context) {
	if a.upd.phase != updateAvailable {
		if circleButton(c, lucide.Settings, locales.SettingsTitle(), 36, false).Clicked() {
			a.showSettings()
		}
		return
	}
	ui.Box(c).Shrink(0).Children(func() {
		if circleButton(c, lucide.Settings, locales.SettingsUpdate(), 36, false).Clicked() {
			a.showSettings()
		}
		ui.Box(c).Size(10, 10).Radius(5).Background(accent).Border(2, panelTop).Absolute().Top(0).Right(0)
	})
}

func areaHeading(c *ui.Context, name string) ui.Element {
	return components.Text(c, name).FontSize(13).Bold().TextColor(inkMuted).SingleLine().Ellipsis("…").Margin(0, 4)
}

// setupView stands in for the tiles until Home Assistant lists the devices:
// it waits, asks for the connection or says what went wrong. top keeps it
// clear of the header.
func (a *application) setupView(c *ui.Context, top float32) {
	ui.Column(c).Fill().Center().Gap(14).Padding(top, 24, 48, 24).Children(func() {
		components.Icon(c, lucide.House).Size(48, 48).TextColor(accent)
		if a.configured() && a.errorKey == "" {
			components.Spinner(c, components.SpinnerRing)
			return
		}
		components.Text(c, locales.SetupTitle()).FontSize(18).Bold().TextColor(ink)
		msg := locales.SetupBody()
		if e := a.errorDisplay(); e != "" {
			msg = e
		}
		components.Text(c, msg).Wrap().TextAlign(ui.Center).FontSize(13).TextColor(inkMuted).MaxWidth(260)
		components.Button(c, locales.Settings(), components.ButtonPrimary).OnClick(a.showSettings)
	})
}

// tile shows one device: what it is, its name and its main state, lit in
// its color while it is on, open or at work. A click switches a device that
// switches on and off, or opens and closes, whose page then opens from the
// button in the corner; it opens the page of any other, an air
// conditioner's controls.
//
// Dragging the tile moves it within its area, sec: into the place of the
// tile under the pointer. Its context menu pins it above the areas, or
// returns it to its area.
func (a *application) tile(c *ui.Context, t ha.Tile, sec ha.Section) {
	s := a.held(t.Entity)
	v := a.tileView(t, s)
	name := tileName(t.Name, sec.Area)
	_, switches := s.Toggle()
	box := ui.Box(c.Key(s.EntityID)).Height(tileHeight)
	if len(sec.Tiles) > 1 {
		d := homeDrag{area: sec.ID, entity: s.EntityID}
		box.Drag(d)
		if box.Dragging() {
			a.dragging(d)
			box.Opacity(0.4)
		}
	}
	if r := a.reorder; r != nil && !r.stale && r.drag.area == sec.ID && r.drag.entity != "" && r.drag.entity != s.EntityID {
		if _, _, over := box.PointerPosition(); over {
			r.move = func() { a.moveTile(sec.ID, r.drag.entity, s.EntityID) }
		}
	}
	box.Children(func() {
		e := ui.ButtonBase(c).Fill().Radius(20).Padding(12).Column().AlignItems(ui.Start)
		if switches {
			e.Label(locales.QuickSwitch(name)).Disabled(!s.Available())
			e.Role(ui.RoleSwitch).Checked(v.lit).Value(v.status)
		} else {
			e.Label(name + ", " + v.status)
		}
		level := float32(0)
		if v.lit {
			level = 1
		}
		level = e.Animate("lit", level, 220*time.Millisecond)
		// The tile is glass, tinted in its color while lit, which swells a
		// little as it is pressed.
		tint := ui.RGBA(255, 255, 255, 0.04).Mix(v.color.Alpha(0.3), level)
		e.Material(glass.Glass{Tint: tint, Interactive: true}).
			Shadow(0, 8, 24, 0, v.color.Alpha(0.28*level))
		if !s.Available() {
			e.Opacity(0.45)
		}
		pin := locales.Pin()
		if sec.ID == pinnedID {
			pin = locales.Unpin()
		}
		e.ContextMenu(func(m *ui.Menu) {
			if m.Item(pin).Chosen() {
				a.togglePin(s.EntityID)
			}
		})
		e.Children(func() {
			deviceDisc(c, v, level, 36)
			ui.Spacer(c)
			components.Text(c, name).FontSize(13).FontWeight(600).TextColor(ink).MaxLines(2)
			components.Text(c, v.status).FontSize(12).TextColor(inkMuted.Mix(v.color.Mix(ink, 0.4), level)).
				SingleLine().Ellipsis("…").Margin(2, 0, 0, 0)
		})
		switch {
		case !e.Clicked():
		case switches:
			a.toggle(s)
		default:
			a.open(t)
		}
		if switches {
			if circleButton(c, lucide.ChevronRight, locales.DeviceDetails(name), 28, false).Absolute().Top(16).Right(12).Clicked() {
				a.open(t)
			}
		}
	})
}

// deviceDisc is the round icon of a device, filled with its color while it
// is lit by level, size DIPs across.
func deviceDisc(c *ui.Context, v tileView, level, size float32) ui.Element {
	disc, glyph := ui.RGBA(255, 255, 255, 0.07).Mix(v.color, level), inkMuted.Mix(discInk, level)
	if v.neutral {
		disc, glyph = v.color.Alpha(0.16), v.color
	}
	icon := size * 19 / 36
	return ui.Box(c).Size(size, size).Radius(size / 2).Shrink(0).Background(disc).Center().Children(func() {
		components.Icon(c, v.icon).Size(icon, icon).TextColor(glyph)
	})
}

// quickSwitch switches the device of entity s, as held, on and off, or
// opens and closes it, if it does either; its element is empty otherwise.
func (a *application) quickSwitch(c *ui.Context, s ha.State, v tileView, name string) ui.Element {
	if _, ok := s.Toggle(); !ok {
		return ui.Element{}
	}
	on := v.lit
	sw := components.Switch(c, &on, locales.QuickSwitch(name), components.SwitchOptions{HideLabel: true})
	sw.Disabled(!s.Available())
	if on != v.lit {
		a.toggle(s)
	}
	return sw
}

// open shows the page of the device of tile t: an air conditioner's
// controls, else what its page shows of any other.
func (a *application) open(t ha.Tile) {
	if t.Entity.Domain() == "climate" {
		a.openClimate(ha.ParseClimate(a.held(t.Entity)))
		return
	}
	a.openDevice(t.Entity.EntityID)
}

// tileName drops the area's name from the front of a device's, which the
// heading above already says: "Lamp" for "Bedroom Lamp" in the bedroom.
func tileName(name, area string) string {
	if area == "" || len(name) <= len(area) || !strings.EqualFold(name[:len(area)], area) {
		return name
	}
	if rest := strings.TrimLeft(name[len(area):], " -_·:："); rest != "" {
		return rest
	}
	return name
}
