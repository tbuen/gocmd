package panel

import (
	"container/list"
	"github.com/tbuen/gocmd/internal/backend/gui"
	"github.com/tbuen/gocmd/internal/backend/tab"
	//"github.com/tbuen/gocmd/internal/config"
	//. "github.com/tbuen/gocmd/internal/global"
	//"github.com/tbuen/gocmd/internal/log"
	"path/filepath"
)

const (
	LEFT = iota
	RIGHT
)

type Header struct {
	Titles []string
	Active int
	// TODO replace by something like guiData:
	// or save this directly in gui
	Offset float64
}

type Panel struct {
	tabs   list.List
	active *list.Element
	header Header
}

var (
	panels [2]Panel
	active *Panel = &panels[LEFT]
)

func Load() {
	// TODO put this in a separate file
	//config.ReadApps()
	//tabcfg, err := config.ReadTabs()
	//if err != nil {
	panels[LEFT].New()
	panels[RIGHT].New()
	//insertTab(PANEL_LEFT, newDefaultDirectory())
	//insertTab(PANEL_RIGHT, newDefaultDirectory())
	//return
	//}
	//if tabcfg.Active < len(panels) {
	//		active = tabcfg.Active
	//	}
	//	for idx := PANEL_LEFT; idx <= PANEL_RIGHT; idx++ {
	//		for _, tab := range tabcfg.Panels[idx].Tabs {
	//			dir := newDirectory(tab.Path, tab.SortKey, tab.SortOrder, tab.Hidden)
	//			insertTab(idx, dir)
	//		}
	//		if tabcfg.Panels[idx].Active < len(panels[idx].tabs) {
	//			panels[idx].active = tabcfg.Panels[idx].Active
	//		} else {
	//			panels[idx].active = 0
	//		}
	//	}
}

func Save() {
	//	tabcfg := config.TabConfig{}
	//	tabcfg.Active = active
	//	for idx := PANEL_LEFT; idx <= PANEL_RIGHT; idx++ {
	//		panel := config.Panel{}
	//		for _, tab := range panels[idx].tabs {
	//			sortKey, sortOrder := tab.dir.SortKey()
	//			panel.Tabs = append(panel.Tabs, config.Tab{tab.dir.Path(), sortKey, sortOrder, tab.dir.Hidden()})
	//		}
	//		panel.Active = panels[idx].active
	//		tabcfg.Panels = append(tabcfg.Panels, panel)
	//	}
	//	config.WriteTabs(&tabcfg)
	// TODO separate file:
}

func Active() *Panel {
	return active
}

func Left() *Panel {
	return &panels[LEFT]
}

func Right() *Panel {
	return &panels[RIGHT]
}

func Toggle() {
	if active == &panels[LEFT] {
		active = &panels[RIGHT]
	} else {
		active = &panels[LEFT]
	}
	gui.Refresh()
}

func (p *Panel) IsActive() bool {
	return p == active
}

func (p *Panel) Header() *Header {
	p.header.Titles = make([]string, p.tabs.Len())
	for e, i := p.tabs.Front(), 0; e != nil; e, i = e.Next(), i+1 {
		// TODO: Directory (or its interface) should become a method Header()
		p.header.Titles[i] = filepath.Base(e.Value.(*tab.Tab).Directory().Path())
		if e == p.active {
			p.header.Active = i
		}
	}
	return &p.header
}

func (p *Panel) First() {
	p.active = p.tabs.Front()
	gui.Refresh()
}

func (p *Panel) Last() {
	p.active = p.tabs.Back()
	gui.Refresh()
}

func (p *Panel) Prev() {
	if e := p.active.Prev(); e != nil {
		p.active = e
		gui.Refresh()
	}
}

func (p *Panel) Next() {
	if e := p.active.Next(); e != nil {
		p.active = e
		gui.Refresh()
	}
}

func (p *Panel) Tab() *tab.Tab {
	return p.active.Value.(*tab.Tab)
}

func (p *Panel) New() {
	p.insert(tab.New())
}

func (p *Panel) Clone() {
	p.insert(p.Tab().Clone())
}

func (p *Panel) Delete() {
	d := p.active
	if p.active.Prev() != nil {
		p.active = p.active.Prev()
	} else {
		p.active = p.active.Next()
	}
	t := p.tabs.Remove(d)
	t.(*tab.Tab).Directory().Destroy()
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
	t.Directory().Reload()
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
