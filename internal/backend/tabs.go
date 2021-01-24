package backend

import (
	"github.com/tbuen/gocmd/internal/config"
	"github.com/tbuen/gocmd/internal/log"
	"path/filepath"
)

const (
	PANEL_LEFT = iota
	PANEL_RIGHT
	PANEL_ACTIVE
)

const (
	TAB_MODE_DIRECTORY = iota
	TAB_MODE_BOOKMARKS
)

type tab struct {
	mode      int
	dir       *Directory
	bookmarks *Bookmarks
}

type panel struct {
	tabs   []tab
	active int
	offset float64
}

type Tabs struct {
	Panel  int
	Titles []string
	Active int
	Offset float64
}

var (
	panels [2]panel
	active int
)

func Load() {
	// TODO put this in a separate file
	config.ReadApps()
	config.ReadBookmarks()
	tabcfg, err := config.ReadTabs()
	if err != nil {
		insertTab(PANEL_LEFT, newDefaultDirectory())
		insertTab(PANEL_RIGHT, newDefaultDirectory())
		return
	}
	if tabcfg.Active < len(panels) {
		active = tabcfg.Active
	}
	for idx := PANEL_LEFT; idx <= PANEL_RIGHT; idx++ {
		for _, tab := range tabcfg.Panels[idx].Tabs {
			dir := newDirectory(tab.Path, tab.SortKey, tab.SortOrder, tab.Hidden)
			insertTab(idx, dir)
		}
		if tabcfg.Panels[idx].Active < len(panels[idx].tabs) {
			panels[idx].active = tabcfg.Panels[idx].Active
		} else {
			panels[idx].active = 0
		}
	}
}

func Save() {
	tabcfg := config.TabConfig{}
	tabcfg.Active = active
	for idx := PANEL_LEFT; idx <= PANEL_RIGHT; idx++ {
		panel := config.Panel{}
		for _, tab := range panels[idx].tabs {
			sortKey, sortOrder := tab.dir.SortKey()
			panel.Tabs = append(panel.Tabs, config.Tab{tab.dir.Path(), sortKey, sortOrder, tab.dir.Hidden()})
		}
		panel.Active = panels[idx].active
		tabcfg.Panels = append(tabcfg.Panels, panel)
	}
	config.WriteTabs(&tabcfg)
	// TODO separate file:
	config.WriteBookmarks()
}

func ActivePanel() int {
	return active
}

func TogglePanel() {
	if active == PANEL_LEFT {
		active = PANEL_RIGHT
	} else {
		active = PANEL_LEFT
	}
	guiRefresh()
}

func GetTabMode(panel int) (mode int) {
	idx := panelIdx(panel)
	mode = panels[idx].tabs[panels[idx].active].mode
	return
}

func ShowBookmarks(panel int) {
	idx := panelIdx(panel)
	if panels[idx].tabs[panels[idx].active].mode == TAB_MODE_DIRECTORY {
		panels[idx].tabs[panels[idx].active].mode = TAB_MODE_BOOKMARKS
		if panels[idx].tabs[panels[idx].active].bookmarks == nil {
			panels[idx].tabs[panels[idx].active].bookmarks = newBookmarks()
		}
		guiRefresh()
	}
}

func HideBookmarks(panel int) {
	idx := panelIdx(panel)
	if panels[idx].tabs[panels[idx].active].mode == TAB_MODE_BOOKMARKS {
		panels[idx].tabs[panels[idx].active].mode = TAB_MODE_DIRECTORY
		guiRefresh()
	}
}

func GetDirectory(panel int) (dir *Directory) {
	idx := panelIdx(panel)
	tab := panels[idx].tabs[panels[idx].active]
	if tab.mode == TAB_MODE_DIRECTORY {
		dir = tab.dir
	}
	return
}

func GetBookmarks(panel int) (bookmarks *Bookmarks) {
	idx := panelIdx(panel)
	tab := panels[idx].tabs[panels[idx].active]
	if tab.mode == TAB_MODE_BOOKMARKS {
		bookmarks = tab.bookmarks
	}
	return
}

func GetTabs(panel int) (tabs *Tabs) {
	tabs = new(Tabs)
	idx := panelIdx(panel)
	tabs.Panel = idx
	for _, tab := range panels[idx].tabs {
		tabs.Titles = append(tabs.Titles, filepath.Base(tab.dir.Path()))
	}
	tabs.Active = panels[idx].active
	tabs.Offset = panels[idx].offset
	return
}

func SetTabOffset(panel int, offset float64) {
	idx := panelIdx(panel)
	panels[idx].offset = offset
}

func CreateTab(panel int) {
	insertTab(panel, newDefaultDirectory())
}

func CloneTab(panel int) {
	idx := panelIdx(panel)
	src := panels[idx].tabs[panels[idx].active].dir
	sortKey, sortOrder := src.SortKey()
	clone := newDirectory(src.Path(), sortKey, sortOrder, src.Hidden())
	insertTab(panel, clone)
}

func DeleteTab(panel int) {
	idx := panelIdx(panel)
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
	guiRefresh()
}

func FirstTab(panel int) {
	idx := panelIdx(panel)
	if panels[idx].active != 0 {
		panels[idx].active = 0
		guiRefresh()
	}
}

func PrevTab(panel int) {
	idx := panelIdx(panel)
	if panels[idx].active > 0 {
		panels[idx].active--
		guiRefresh()
	}
}

func NextTab(panel int) {
	idx := panelIdx(panel)
	num := len(panels[idx].tabs)
	if panels[idx].active < num-1 {
		panels[idx].active++
		guiRefresh()
	}
}

func LastTab(panel int) {
	idx := panelIdx(panel)
	num := len(panels[idx].tabs)
	if panels[idx].active != num-1 {
		panels[idx].active = num - 1
		guiRefresh()
	}
}

func insertTab(panel int, dir *Directory) {
	idx := panelIdx(panel)
	tabs := &panels[idx].tabs
	active := &panels[idx].active
	log.Println(log.TAB, "creating tab, before:", len(*tabs))
	if len(*tabs) == 0 {
		tab := tab{TAB_MODE_DIRECTORY, dir, nil}
		*tabs = append(*tabs, tab)
		*active = 0
	} else {
		*tabs = append(*tabs, tab{})
		copy((*tabs)[*active+2:], (*tabs)[*active+1:])
		tab := tab{TAB_MODE_DIRECTORY, dir, nil}
		(*tabs)[*active+1] = tab
		*active++
	}
	log.Println(log.TAB, "creating tab, after:", len(*tabs))
	dir.Reload()
}

func panelIdx(panel int) (idx int) {
	if panel == PANEL_ACTIVE {
		idx = active
	} else {
		idx = panel
	}
	return
}
