package fs

const (
	TAB_LEFT  = 0
	TAB_RIGHT = 1
)

var tabs [2][]*dir

func init() {
	left := &dir{"/home/thomas"}
	right := &dir{"/home/thomas/github/gocmd"}
	tabs[TAB_LEFT] = append(tabs[TAB_LEFT], left)
	tabs[TAB_RIGHT] = append(tabs[TAB_RIGHT], right)
}

func Tab(index int) *dir {
	return tabs[index][0]
}
