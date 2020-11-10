package backend

import (
	"github.com/tbuen/gocmd/internal/log"
)

const (
	PANEL_LEFT = iota
	PANEL_RIGHT
	PANEL_ACTIVE
)

type panel struct {
	tabs   []Directory
	active int
}

var (
	panels [2]panel
	active int
)

func Load() {
	CreateTab(PANEL_LEFT)
	CreateTab(PANEL_RIGHT)
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

func GetDirectory(panel int) (d Directory) {
	idx := panelIdx(panel)
	if len(panels[idx].tabs) > 0 {
		d = panels[idx].tabs[panels[idx].active]
	}
	return
}

func Tabs(panel int) (titles []string, active int) {
	idx := panelIdx(panel)
	for _, d := range panels[idx].tabs {
		titles = append(titles, d.Path())
	}
	active = panels[idx].active
	return
}

func CreateTab(panel int) {
	insertTab(panel, newDirectory(""))
}

func DuplicateTab(panel int) {
	org := GetDirectory(panel)
	if org != nil {
		dupl := newDirectory(org.Path())
		dupl.SetSort(org.Sort())
		insertTab(panel, dupl)
	}
}

func DeleteTab(panel int) {
	idx := panelIdx(panel)
	tabs := &panels[idx].tabs
	active := &panels[idx].active
	d := (*tabs)[*active]
	log.Println(log.TAB, "deleting tab, before:", len(*tabs))
	*tabs = append((*tabs)[:*active], (*tabs)[*active+1:]...)
	if len(*tabs) == 0 {
		CreateTab(panel)
	}
	if *active >= len(*tabs) {
		*active = len(*tabs) - 1
	}
	log.Println(log.TAB, "deleting tab, after:", len(*tabs))
	d.Destroy()
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

func insertTab(panel int, d Directory) {
	idx := panelIdx(panel)
	tabs := &panels[idx].tabs
	active := &panels[idx].active
	log.Println(log.TAB, "creating tab, before:", len(*tabs))
	if len(*tabs) == 0 {
		*tabs = append(*tabs, d)
		*active = 0
	} else {
		*tabs = append(*tabs, nil)
		copy((*tabs)[*active+2:], (*tabs)[*active+1:])
		(*tabs)[*active+1] = d
		*active++
	}
	log.Println(log.TAB, "creating tab, after:", len(*tabs))
	d.Reload()
}

func panelIdx(panel int) (idx int) {
	if panel == PANEL_ACTIVE {
		idx = active
	} else {
		idx = panel
	}
	return
}
