package backend

var refreshFct func()

func RegisterRefresh(refresh func()) {
	refreshFct = refresh
}

func guiRefresh() {
	if refreshFct != nil {
		refreshFct()
	}
}
