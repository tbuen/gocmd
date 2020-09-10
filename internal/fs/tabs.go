package fs

const (
	TAB_LEFT  = 0
	TAB_RIGHT = 1
)

var tabs [2][]Directory

func init() {
	left := newDirectory("/home/thomas")
	right := newDirectory("/home/thomas/github/gocmd")
	tabs[TAB_LEFT] = append(tabs[TAB_LEFT], left)
	tabs[TAB_RIGHT] = append(tabs[TAB_RIGHT], right)
	tabs[TAB_LEFT][0].reload()
	tabs[TAB_RIGHT][0].reload()
}

func Tab(index int) Directory {
	return tabs[index][0]
}
