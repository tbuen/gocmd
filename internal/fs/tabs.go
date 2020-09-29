package fs

const (
	PANEL_LEFT = iota
	PANEL_RIGHT
	PANEL_ACTIVE
)

var (
	tabs   [2][]Directory
	active = PANEL_LEFT
)

func Load() {
	left := newDirectory("")
	right := newDirectory("")
	tabs[PANEL_LEFT] = append(tabs[PANEL_LEFT], left)
	tabs[PANEL_RIGHT] = append(tabs[PANEL_RIGHT], right)
	tabs[PANEL_LEFT][0].Reload()
	tabs[PANEL_RIGHT][0].Reload()
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
	switch panel {
	case PANEL_LEFT, PANEL_RIGHT:
		if len(tabs[panel]) > 0 {
			d = tabs[panel][0]
		}
	case PANEL_ACTIVE:
		if len(tabs[active]) > 0 {
			d = tabs[active][0]
		}
	}
	return
}
