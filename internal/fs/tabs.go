package fs

const (
	TAB_LEFT   = 0
	TAB_RIGHT  = 1
	TAB_ACTIVE = 2
)

var (
	tabs   [2][]Directory
	active int = 1
)

func init() {
	left := newDirectory("/home/thomas")
	right := newDirectory("/home/thomas/github/gocmd")
	tabs[TAB_LEFT] = append(tabs[TAB_LEFT], left)
	tabs[TAB_RIGHT] = append(tabs[TAB_RIGHT], right)
	tabs[TAB_LEFT][0].Reload()
	tabs[TAB_RIGHT][0].Reload()
}

func GetDirectory(index int) Directory {
	switch index {
	case TAB_LEFT, TAB_RIGHT:
		return tabs[index][0]
	case TAB_ACTIVE:
		return tabs[active][0]
	default:
		return nil
	}
}
