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
	left = iota
	right
)

type Header struct {
	Titles []string
	Active int
	Offset float64
}

type Panel struct {
	tabs   list.List
	active *list.Element
	header Header
}

var (
	panels [2]Panel
	active *Panel = &panels[left]
)

func Load() {
	// TODO put this in a separate file
	//config.ReadApps()
	//tabcfg, err := config.ReadTabs()
	//if err != nil {
	panels[left].NewTab()
	panels[right].NewTab()
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
	return &panels[left]
}

func Right() *Panel {
	return &panels[right]
}

func Toggle() {
	if active == &panels[left] {
		active = &panels[right]
	} else {
		active = &panels[left]
	}
	gui.Refresh()
}

func (p *Panel) IsActive() bool {
	return p == active
}

func (p *Panel) Header() *Header {
	p.header.Titles = make([]string, 0, p.tabs.Len())
	for e := p.tabs.Front(); e != nil; e = e.Next() {
		p.header.Titles = append(p.header.Titles, filepath.Base(e.Value.(*tab.Tab).Directory().Path()))
	}
	p.header.Active = 0
	return &p.header
}

func (p *Panel) Tab() *tab.Tab {
	return p.active.Value.(*tab.Tab)
}

func (p *Panel) NewTab() {
	p.insert(tab.New())
}

func (p *Panel) CloneTab() {
	p.insert(p.Tab().Clone())
}

func (p *Panel) DeleteTab() {
	/*idx := panelIdx(panel)
	tabs := &panels[idx].tabs
	active := &panels[idx].active
	dir := (*tabs)[*active].dir
	dir.Destroy()
	log.Println(log.TAB, "deleting tab, before:", len(*tabs))
	*tabs = append((*tabs)[:*active], (*tabs)[*active+1:]...)
	if len(*tabs) == 0 {
		CreateTab(panel)
	}
	if *active > 0 {
		*active--
	}
	log.Println(log.TAB, "deleting tab, after:", len(*tabs))
	gui.Refresh()*/
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
}

func FirstTab(panel int) {
	idx := panelIdx(panel)
	if panels[idx].active != 0 {
		panels[idx].active = 0
		gui.Refresh()
	}
}

func PrevTab(panel int) {
	idx := panelIdx(panel)
	if panels[idx].active > 0 {
		panels[idx].active--
		gui.Refresh()
	}
}

func NextTab(panel int) {
	idx := panelIdx(panel)
	num := len(panels[idx].tabs)
	if panels[idx].active < num-1 {
		panels[idx].active++
		gui.Refresh()
	}
}

func LastTab(panel int) {
	idx := panelIdx(panel)
	num := len(panels[idx].tabs)
	if panels[idx].active != num-1 {
		panels[idx].active = num - 1
		gui.Refresh()
	}
}

func panelIdx(panel int) (idx int) {
	if panel == PANEL_ACTIVE {
		idx = active
	} else {
		idx = panel
	}
	return
}
*/
