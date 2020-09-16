package fs

const (
	PANEL_LEFT   = 0
	PANEL_RIGHT  = 1
	PANEL_ACTIVE = 2
)

var (
	tabs   [2][]Directory
	active int = PANEL_LEFT
)

func init() {
	left := newDirectory("/home/thomas")
	right := newDirectory("/home/thomas/github/gocmd")
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

func GetDirectory(index int) Directory {
	switch index {
	case PANEL_LEFT, PANEL_RIGHT:
		return tabs[index][0]
	case PANEL_ACTIVE:
		return tabs[active][0]
	default:
		return nil
	}
}
