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

type panel struct {
	tabs   []*Directory
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
	config.ReadApps()
	tabcfg, err := config.ReadTabs()
	if err != nil {
		insertTab(PANEL_LEFT, newDirectory("", config.SORT_BY_NAME, config.SORT_ASCENDING, false))
		insertTab(PANEL_LEFT, newDirectory("", config.SORT_BY_NAME, config.SORT_ASCENDING, false))
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
		for _, dir := range panels[idx].tabs {
			sortKey, sortOrder := dir.SortKey()
			panel.Tabs = append(panel.Tabs, config.Tab{dir.Path(), sortKey, sortOrder, dir.Hidden()})
		}
		panel.Active = panels[idx].active
		tabcfg.Panels = append(tabcfg.Panels, panel)
	}
	config.WriteTabs(&tabcfg)
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

func GetDirectory(panel int) (dir *Directory) {
	idx := panelIdx(panel)
	if len(panels[idx].tabs) > 0 {
		dir = panels[idx].tabs[panels[idx].active]
	}
	return
}

func GetTabs(panel int) (tabs *Tabs) {
	tabs = new(Tabs)
	idx := panelIdx(panel)
	tabs.Panel = idx
	for _, dir := range panels[idx].tabs {
		tabs.Titles = append(tabs.Titles, filepath.Base(dir.Path()))
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
	insertTab(panel, newDirectory("", config.SORT_BY_NAME, config.SORT_ASCENDING, false))
}

func CloneTab(panel int) {
	src := GetDirectory(panel)
	if src != nil {
		sortKey, sortOrder := src.SortKey()
		clone := newDirectory(src.Path(), sortKey, sortOrder, src.Hidden())
		insertTab(panel, clone)
	}
}

func DeleteTab(panel int) {
	idx := panelIdx(panel)
	tabs := &panels[idx].tabs
	active := &panels[idx].active
	dir := (*tabs)[*active]
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
		*tabs = append(*tabs, dir)
		*active = 0
	} else {
		*tabs = append(*tabs, nil)
		copy((*tabs)[*active+2:], (*tabs)[*active+1:])
		(*tabs)[*active+1] = dir
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
