package gui

var refreshFct func()

func RegisterRefresh(refresh func()) {
	refreshFct = refresh
}

func Refresh() {
	if refreshFct != nil {
		refreshFct()
	}
}
