package backend

import (
	"github.com/tbuen/gocmd/internal/backend/panel"
)

func Start() {
	panel.Load()
}

func Stop() {
	panel.Save()
}
