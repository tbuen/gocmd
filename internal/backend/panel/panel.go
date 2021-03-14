// Package panel implements the two panels and functions to access their tabs.
package panel

import (
	"container/list"
	"path/filepath"

	"github.com/tbuen/gocmd/internal/backend/gui"
	"github.com/tbuen/gocmd/internal/backend/tab"
	"github.com/tbuen/gocmd/internal/config"
	//. "github.com/tbuen/gocmd/internal/global"
	//"github.com/tbuen/gocmd/internal/log"
)

const (
	left = iota
	right
)

// TODO am besten ganze header struktur hier raus und zugriff auf alle tabs erlauben,
// die dann eine Header-Methode haben
type Header struct {
	Titles []string
	Active int
	// TODO replace by something like guiData:
	// or save this directly in gui
	Offset float64
}

// Panel is one of the two panels.
type Panel struct {
	tabs   list.List
	active *list.Element
	header Header
}

var (
	panels [2]Panel
	active *Panel = &panels[left]
)

// Load loads the tabs of both panels from the configuration file.
func Load() {
	c := config.GetTabs()
	for i := left; i <= right; i++ {
		if i < len(c.Panels) && len(c.Panels[i].Tabs) > 0 {
			for j, ct := range c.Panels[i].Tabs {
				t := tab.NewWithConfig(ct)
				e := panels[i].tabs.PushBack(t)
				if j == c.Panels[i].Active {
					panels[i].active = e
				}
				t.Reload()
			}
		} else {
			panels[i].New()
		}
		if i == c.Active {
			active = &panels[i]
		}
	}
}

// Save saves the tabs of both panels to the configuration file.
func Save() {
	c := config.GetTabs()
	c.Panels = c.Panels[:0]
	for i := left; i <= right; i++ {
		cp := config.Panel{}
		for j, e := 0, panels[i].tabs.Front(); e != nil; j, e = j+1, e.Next() {
			cp.Tabs = append(cp.Tabs, e.Value.(*tab.Tab).Config())
			if e == panels[i].active {
				cp.Active = j
			}
		}
		c.Panels = append(c.Panels, cp)
		if &panels[i] == active {
			c.Active = i
		}
	}
}

// Active returns the active panel.
func Active() *Panel {
	return active
}

// Left returns the left panel.
func Left() *Panel {
	return &panels[left]
}

// Right returns the right panel.
func Right() *Panel {
	return &panels[right]
}

// Toggle toggles the active panel.
func Toggle() {
	if active == &panels[left] {
		active = &panels[right]
	} else {
		active = &panels[left]
	}
	gui.Refresh()
}

// IsActive returns if the panel is active.
func (p *Panel) IsActive() bool {
	return p == active
}

func (p *Panel) Header() *Header {
	p.header.Titles = make([]string, p.tabs.Len())
	for e, i := p.tabs.Front(), 0; e != nil; e, i = e.Next(), i+1 {
		// TODO: Directory (or its interface) or even better the tab(!) should become a method Header()
		p.header.Titles[i] = filepath.Base(e.Value.(*tab.Tab).Directory().Path())
		if e == p.active {
			p.header.Active = i
		}
	}
	return &p.header
}

// First activates the first tab of the panel.
func (p *Panel) First() {
	p.active = p.tabs.Front()
	gui.Refresh()
}

// Last activates the last tab of the panel.
func (p *Panel) Last() {
	p.active = p.tabs.Back()
	gui.Refresh()
}

// Prev activates the previous tab of the panel if available.
func (p *Panel) Prev() {
	e := p.active.Prev()
	if e == nil {
		return
	}
	p.active = e
	gui.Refresh()
}

// Next activates the next tab of the panel if available.
func (p *Panel) Next() {
	e := p.active.Next()
	if e == nil {
		return
	}
	p.active = e
	gui.Refresh()
}

// Tab returns the active tab of the panel.
func (p *Panel) Tab() *tab.Tab {
	return p.active.Value.(*tab.Tab)
}

// New creates a new tab in the panel.
func (p *Panel) New() {
	p.insert(tab.New())
}

// Clone duplicates the active tab of the panel.
func (p *Panel) Clone() {
	p.insert(p.Tab().Clone())
}

// Delete deletes the active tab of the panel.
func (p *Panel) Delete() {
	d := p.active
	if p.active.Prev() != nil {
		p.active = p.active.Prev()
	} else {
		p.active = p.active.Next()
	}
	t := p.tabs.Remove(d)
	t.(*tab.Tab).Destroy()
	if p.tabs.Len() == 0 {
		p.New()
	} else {
		gui.Refresh()
	}
}

func (p *Panel) insert(t *tab.Tab) {
	if p.active == nil {
		p.active = p.tabs.PushFront(t)
	} else {
		p.active = p.tabs.InsertAfter(t, p.active)
	}
	t.Reload()
}

/*func AddBookmark(panel int) {
	idx := panelIdx(panel)
	if panels[idx].tabs[panels[idx].active].mode == TAB_MODE_DIRECTORY {
		config.Bookmarks().Add(Bookmark{panels[idx].tabs[panels[idx].active].dir.Path()})
		gui.Refresh()
	}
}*/

/*func SetTabOffset(panel int, offset float64) {
	idx := panelIdx(panel)
	panels[idx].offset = offset
}*/
